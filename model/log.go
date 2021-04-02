package model

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var SugarLogger *zap.SugaredLogger

func InitLogger() {
	writerSyncer := getLogWriter() // 콘솔창 출력 & 파일쓰기
	encoder := getEncoder()
	consoleErrors := zapcore.Lock(os.Stderr) // 콘솔창에만 출력

	// zapcore.NewTee => 여러 로그레벨(core 인터페이스)을 반환하려고 할 경우 사용.
	coreError := zapcore.NewTee(
		 //zapcore.NewCore(encoder, writerSyncer, zapcore.ErrorLevel),
		 zapcore.NewCore(encoder, consoleErrors, zapcore.InfoLevel),
		 zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel), // ErrorLevel, DebugLevel 만 파일에 기록한다.
		 //zapcore.NewCore(encoder, consoleErrors, zapcore.DebugLevel),
		)
	errorLog := zap.New(coreError, zap.AddCaller())

	SugarLogger = errorLog.Sugar()
}

// 로그 인코딩 형식 커스터마이징
func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig() // 배포된 이후에 사용하는 설정
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // timestamp ms 단위를 ISO8601 형식으로 변환
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 대문자로 직렬화 (info -> INFO)
	return zapcore.NewConsoleEncoder(encoderConfig)

	// JSON 형식으로 저장할 경우.
	// return zapcore.NewJSONEncoder(encoderConfig)
}

// 로그를 파일로 기록할 메서드
// 로그파일 로테이션 라이브러리 lumberjack 으로 로그파일 보관조건 입력
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: "./logfile.log",
		MaxSize: 80, // 80MB 까지 보관, 최대 100MB 가능
		MaxBackups: 5, // 최대 백업파일 개수 (5개 넘어가면 삭제)
		MaxAge: 30, // 30일 까지 보관
		Compress: false,
	}
	//file, _ := os.Create("./logfile.log")
	return zapcore.AddSync(lumberJackLogger)
}