package services

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParsedSignal struct {
	Signal     []float64
	SampleRate float64
	FileName   string
}

func ParseSignalFile(filename string, content []byte) (*ParsedSignal, error) {
	lowerName := strings.ToLower(filename)

	if strings.HasSuffix(lowerName, ".json") {
		return parseJSON(content, filename)
	}
	if strings.HasSuffix(lowerName, ".csv") || strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".dat") {
		return parseCSV(content, filename)
	}

	return nil, fmt.Errorf("不支持的文件格式: %s，仅支持 .csv、.json、.txt、.dat", filename)
}

func parseJSON(content []byte, filename string) (*ParsedSignal, error) {
	var raw interface{}
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	result := &ParsedSignal{FileName: filename, SampleRate: 100.0}

	switch v := raw.(type) {
	case []interface{}:
		sig, err := toFloatSlice(v)
		if err != nil {
			return nil, err
		}
		result.Signal = sig
	case map[string]interface{}:
		if sigRaw, ok := v["signal"]; ok {
			if arr, ok := sigRaw.([]interface{}); ok {
				sig, err := toFloatSlice(arr)
				if err != nil {
					return nil, err
				}
				result.Signal = sig
			} else {
				return nil, fmt.Errorf("JSON 中 'signal' 字段不是数组")
			}
		} else {
			return nil, fmt.Errorf("JSON 中缺少 'signal' 字段")
		}
		if sr, ok := v["sampleRate"].(float64); ok {
			result.SampleRate = sr
		} else if sr, ok := v["sample_rate"].(float64); ok {
			result.SampleRate = sr
		}
	default:
		return nil, fmt.Errorf("JSON 格式不正确，需要是数组或包含 signal 字段的对象")
	}

	if len(result.Signal) == 0 {
		return nil, fmt.Errorf("未找到有效的信号数据")
	}

	return result, nil
}

func parseCSV(content []byte, filename string) (*ParsedSignal, error) {
	text := string(content)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("文件为空")
	}

	signal := make([]float64, 0, len(lines))
	startLine := 0

	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if strings.HasPrefix(firstLine, "#") || strings.HasPrefix(firstLine, "//") {
			startLine = 1
		} else {
			if _, err := strconv.ParseFloat(strings.Split(firstLine, ",")[0], 64); err != nil {
				startLine = 1
			}
		}
	}

	for i := startLine; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ',' || r == ';' || r == '\t' || r == ' '
		})
		if len(fields) == 0 {
			continue
		}

		lastField := strings.TrimSpace(fields[len(fields)-1])
		val, err := strconv.ParseFloat(lastField, 64)
		if err != nil {
			continue
		}
		signal = append(signal, val)
	}

	if len(signal) == 0 {
		return nil, fmt.Errorf("未找到有效的数值信号数据")
	}

	return &ParsedSignal{
		Signal:     signal,
		SampleRate: 100.0,
		FileName:   filename,
	}, nil
}

func toFloatSlice(arr []interface{}) ([]float64, error) {
	result := make([]float64, len(arr))
	for i, v := range arr {
		switch val := v.(type) {
		case float64:
			result[i] = val
		case float32:
			result[i] = float64(val)
		case int:
			result[i] = float64(val)
		case int64:
			result[i] = float64(val)
		case string:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("第 %d 个元素不是有效数字: %s", i, val)
			}
			result[i] = f
		default:
			return nil, fmt.Errorf("第 %d 个元素类型不支持: %T", i, v)
		}
	}
	return result, nil
}

func ReadAll(r io.Reader, maxSize int64) ([]byte, error) {
	lr := io.LimitReader(r, maxSize)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) >= maxSize {
		return nil, fmt.Errorf("文件过大，超过限制 %d MB", maxSize/1024/1024)
	}
	return data, nil
}
