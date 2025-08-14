package rtype

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInt64Array_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Int64Array
		expected string
		wantErr  bool
	}{
		{
			name:     "empty array",
			input:    Int64Array{},
			expected: "[]",
			wantErr:  false,
		},
		{
			name:     "single element",
			input:    Int64Array{123},
			expected: `["123"]`,
			wantErr:  false,
		},
		{
			name:     "multiple elements",
			input:    Int64Array{123, 456, 789},
			expected: `["123","456","789"]`,
			wantErr:  false,
		},
		{
			name:     "negative numbers",
			input:    Int64Array{-123, 0, 456},
			expected: `["-123","0","456"]`,
			wantErr:  false,
		},
		{
			name:     "large numbers",
			input:    Int64Array{9223372036854775807, -9223372036854775808},
			expected: `["9223372036854775807","-9223372036854775808"]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64Array.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("Int64Array.MarshalJSON() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}

func TestInt64Array_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64Array
		wantErr  bool
	}{
		{
			name:     "empty array",
			input:    "[]",
			expected: Int64Array{},
			wantErr:  false,
		},
		{
			name:     "single element",
			input:    `["123"]`,
			expected: Int64Array{123},
			wantErr:  false,
		},
		{
			name:     "multiple elements",
			input:    `["123","456","789"]`,
			expected: Int64Array{123, 456, 789},
			wantErr:  false,
		},
		{
			name:     "negative numbers",
			input:    `["-123","0","456"]`,
			expected: Int64Array{-123, 0, 456},
			wantErr:  false,
		},
		{
			name:     "with spaces",
			input:    `["123", "456", "789"]`,
			expected: Int64Array{123, 456, 789},
			wantErr:  false,
		},
		{
			name:    "invalid json",
			input:   `["123" "456"]`, // missing comma
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   `["123a","456"]`,
			wantErr: true,
		},
		{
			name:    "not an array",
			input:   `"not an array"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Int64Array
			err := result.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64Array.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Int64Array.UnmarshalJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInt64Array_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input Int64Array
	}{
		{
			name:  "empty array",
			input: Int64Array{},
		},
		{
			name:  "single element",
			input: Int64Array{123},
		},
		{
			name:  "multiple elements",
			input: Int64Array{123, 456, 789},
		},
		{
			name:  "negative numbers",
			input: Int64Array{-123, 0, 456},
		},
		{
			name:  "large numbers",
			input: Int64Array{9223372036854775807, -9223372036854775808},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 序列化
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			// 反序列化
			var result Int64Array
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// 验证结果
			if !reflect.DeepEqual(result, tt.input) {
				t.Errorf("Round trip failed, got %v, want %v", result, tt.input)
			}
		})
	}
}

func TestInt64Array_UnmarshalJSON_PointerReceiver(t *testing.T) {
	var arr *Int64Array
	data := []byte(`["123","456"]`)

	// 测试 nil 指针的 Unmarshal
	err := arr.UnmarshalJSON(data)
	if err == nil {
		t.Error("Expected error when unmarshaling to nil pointer, got nil")
	}

	// 测试正常指针的 Unmarshal
	arr = &Int64Array{}
	err = arr.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("Unexpected error when unmarshaling to valid pointer: %v", err)
	}

	expected := Int64Array{123, 456}
	if !reflect.DeepEqual(*arr, expected) {
		t.Errorf("Expected %v, got %v", expected, *arr)
	}
}



func BenchmarkInt64Array_MarshalJSON(b *testing.B) {
	arr := Int64Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := arr.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInt64Array_UnmarshalJSON(b *testing.B) {
	data := []byte(`["1","2","3","4","5","6","7","8","9","10"]`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var arr Int64Array
		err := arr.UnmarshalJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
