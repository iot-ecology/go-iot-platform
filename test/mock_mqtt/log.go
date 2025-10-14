/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var myTimeEncoder = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 按照 "2006-01-02 15:04:05" 的格式编码时间
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
})

func InitLog() {
	encoderConfig := zapcore.EncoderConfig{
		// 使用自定义的时间编码器
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码日志级别
		EncodeTime:     myTimeEncoder,                 // 使用自定义的时间编码器
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码调用者
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 使用 Console 编码器
		getWriteSync(),
		zap.NewAtomicLevelAt(zap.ErrorLevel),      // 设置日志级别为 Debug

	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger) // 替换全局 Logger
	//确保日志被刷新
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil && !errors.Is(err, syscall.ENOTTY) {
	//		zap.S().Errorf("日志同步失败 %+v", err)
	//	}
	//}(logger)

	// 记录一条日志作为示例
	logger.Debug("这是一个调试级别的日志")
}

func getWriteSync() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./MQTT-TEST.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	syncFile := zapcore.AddSync(lumberJackLogger)
	syncConsole := zapcore.AddSync(os.Stderr)
	return zapcore.NewMultiWriteSyncer(syncConsole, syncFile)
}