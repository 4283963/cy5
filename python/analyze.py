#!/usr/bin/env python3
"""Ultrasonic NDT defect-detection signal processor.

Reads ONE JSON object from STDIN and writes ONE JSON object to STDOUT.
All diagnostic logging is written to STDERR (never STDOUT) so that a
calling service (e.g. a Go backend) can parse stdout as pure JSON.
"""
import json
import math
import sys

import numpy as np
from scipy.ndimage import gaussian_filter1d
from scipy.signal import find_peaks

DEFAULTS = {
    "gaussian_sigma": 3,
    "fft_window": "hann",
    "peak_prominence": 0.2,
    "peak_distance": 10,
}

STEEL_VELOCITY_MM_PER_US = 5.9  # 5900 m/s == 5.9 mm/us
DEFAULT_SAMPLE_RATE = 100.0


def log(message):
    print(message, file=sys.stderr, flush=True)


def empty_result():
    return {
        "filtered": [],
        "fft_freq": [],
        "fft_mag": [],
        "peaks": [],
        "stats": {
            "points": 0,
            "snr_improvement": 0.0,
            "peak_count": 0,
            "max_amplitude": 0.0,
        },
    }


def to_float_list(values):
    return [float(v) for v in np.asarray(values, dtype=float).ravel()]


def resolve_params(raw_params):
    params = dict(DEFAULTS)
    if isinstance(raw_params, dict):
        for key, default in DEFAULTS.items():
            value = raw_params.get(key, default)
            if value is not None:
                params[key] = value
    return params


def build_window(name, n):
    key = str(name or "rect").strip().lower()
    if key == "hann":
        return np.hanning(n)
    if key == "hamming":
        return np.hamming(n)
    return np.ones(n)


def compute_fft(filtered, sample_rate, window_name):
    n = len(filtered)
    if n == 0:
        return [], []
    windowed = filtered * build_window(window_name, n)
    spectrum = np.fft.rfft(windowed)
    freq = np.fft.rfftfreq(n, d=1.0 / float(sample_rate))
    mag = np.abs(spectrum) * 2.0 / n
    return to_float_list(freq), to_float_list(mag)


def detect_peaks(filtered, sample_rate, prominence, distance):
    if len(filtered) == 0:
        return []
    peaks, properties = find_peaks(
        filtered, prominence=float(prominence), distance=int(distance)
    )
    prominences = properties.get("prominences")
    if prominences is None:
        prominences = [0.0] * len(peaks)
    results = []
    for offset, index in enumerate(peaks):
        idx = int(index)
        time_us = idx / float(sample_rate)
        depth_mm = time_us * STEEL_VELOCITY_MM_PER_US / 2.0
        results.append({
            "index": idx,
            "time": float(time_us),
            "depth": float(depth_mm),
            "amplitude": float(filtered[idx]),
            "prominence": float(prominences[offset]),
        })
    return results


def compute_snr(raw, filtered):
    residual = np.asarray(raw, dtype=float) - np.asarray(filtered, dtype=float)
    std_raw = float(np.std(raw))
    std_residual = float(np.std(residual))
    if std_raw <= 0.0 or std_residual <= 0.0:
        return 0.0
    return float(20.0 * math.log10(std_raw / std_residual))


def process(payload):
    signal = payload.get("signal", []) or []
    raw = np.asarray(signal, dtype=float)
    sample_rate = float(payload.get("sample_rate", DEFAULT_SAMPLE_RATE) or DEFAULT_SAMPLE_RATE)
    if sample_rate <= 0.0:
        log("warning: non-positive sample_rate; falling back to default")
        sample_rate = DEFAULT_SAMPLE_RATE

    params = resolve_params(payload.get("params", {}))

    if raw.size == 0:
        return empty_result()

    filtered = gaussian_filter1d(raw, sigma=float(params["gaussian_sigma"]))
    fft_freq, fft_mag = compute_fft(filtered, sample_rate, params["fft_window"])
    peaks = detect_peaks(
        filtered,
        sample_rate,
        float(params["peak_prominence"]),
        int(params["peak_distance"]),
    )
    stats = {
        "points": int(raw.size),
        "snr_improvement": compute_snr(raw, filtered),
        "peak_count": int(len(peaks)),
        "max_amplitude": float(np.max(np.abs(filtered))) if filtered.size else 0.0,
    }
    return {
        "filtered": to_float_list(filtered),
        "fft_freq": fft_freq,
        "fft_mag": fft_mag,
        "peaks": peaks,
        "stats": stats,
    }


def main():
    raw_text = sys.stdin.read()
    try:
        payload = json.loads(raw_text) if raw_text.strip() else {}
        if not isinstance(payload, dict):
            raise ValueError("input must be a JSON object")
    except Exception as exc:
        log("error: failed to parse input JSON: {}".format(exc))
        sys.stdout.write(json.dumps(empty_result()))
        return

    try:
        result = process(payload)
    except Exception as exc:
        import traceback
        log("error: processing failed: {}".format(exc))
        log(traceback.format_exc())
        result = empty_result()

    sys.stdout.write(json.dumps(result))
    sys.stdout.write("\n")


if __name__ == "__main__":
    main()
