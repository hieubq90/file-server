package main

import (
	"file-server/database"
	"file-server/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title File Management Service
// @version 1.0
// @description This is an API for File Server Application

// @contact.name Bùi Quang Hiếu
// @contact.email hieubq90@gmail.com

// @BasePath /api
func main() {
	// Khởi tạo kết nối CSDL
	if err := database.Connect(); err != nil {
		panic(fmt.Errorf("lỗi tải cấu hình: %s \n", err.Error()))
	}

	// Khởi tạo dịch vụ API
	server.InitServer()

	// Chạy API với 1 goroutine khác
	go func() {
		if err := server.StartServer(); err != nil {
			log.Panic(err)
		}
	}()

	// Cài đặt Graceful shutdown
	c := make(chan os.Signal, 1)                    // Tạo channel lắng nghe tín hiệu
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Gửi thông báo cho channel khi có tín hiệu dừng

	_ = <-c // Dừng main thread cho đến khi nhận được tín hiệu dừng

	// Đóng API & dọn dẹp tài nguyên
	fmt.Println("Chuẩn bị dừng dịch vụ...")
	server.Shutdown()

	fmt.Println("Dọn dẹp tài nguyên & ngắt các kết nối...")
	// database.Close()
	fmt.Println("Đã dừng dịch vụ.")

}
