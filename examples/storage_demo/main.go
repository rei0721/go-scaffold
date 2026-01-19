package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	"github.com/rei0721/go-scaffold/pkg/storage"
)

func main() {
	// 创建文件服务实例
	cfg := &storage.Config{
		FSType:   storage.FSTypeOS,
		BasePath: ".",
	}

	fs, err := storage.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	fmt.Println("=== 文件服务工具库示例 ===")

	// 1. 基础文件操作
	fmt.Println("1. 基础文件操作")
	testBasicFileOps(fs)

	// 2. MIME检测
	fmt.Println("\n2. MIME类型检测")
	testMIMEDetection(fs)

	// 3. 文件复制
	fmt.Println("\n3. 文件复制")
	testFileCopy(fs)

	// 4. Excel处理 (示例)
	fmt.Println("\n4. Excel处理")
	testExcelOps(fs)

	// 5. 图片处理 (示例)
	fmt.Println("\n5. 图片处理")
	testImageOps(fs)

	// 6. 文件监听 (示例)
	fmt.Println("\n6. 文件监听")
	testFileWatch(fs)

	fmt.Println("\n=== 示例完成 ===")
}

// 测试基础文件操作
func testBasicFileOps(fs storage.Storage) {
	// 写入文件
	err := fs.WriteFile("test_demo.txt", []byte("Hello, FileService!"), 0644)
	if err != nil {
		fmt.Printf("  写入文件失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ 写入文件成功: test_demo.txt")

	// 读取文件
	data, err := fs.ReadFile("test_demo.txt")
	if err != nil {
		fmt.Printf("  读取文件失败: %v\n", err)
		return
	}
	fmt.Printf("  ✓ 读取文件成功: %s\n", string(data))

	// 检查文件是否存在
	exists, _ := fs.Exists("test_demo.txt")
	fmt.Printf("  ✓ 文件存在检查: %v\n", exists)

	// 获取文件大小
	size, _ := fs.FileSize("test_demo.txt")
	fmt.Printf("  ✓ 文件大小: %d 字节\n", size)

	// 清理
	fs.Remove("test_demo.txt")
}

// 测试MIME检测
func testMIMEDetection(fs storage.Storage) {
	// 从字节检测
	mimeType, err := fs.DetectMIMEFromBytes([]byte("Hello, World!"))
	if err != nil {
		fmt.Printf("  MIME检测失败: %v\n", err)
		return
	}
	fmt.Printf("  ✓ 文本数据MIME: %s\n", mimeType)

	// 创建一个简单的HTML文件
	htmlData := []byte("<!DOCTYPE html><html><body>Test</body></html>")
	fs.WriteFile("test_demo.html", htmlData, 0644)

	mimeType, err = fs.DetectMIME("test_demo.html")
	if err != nil {
		fmt.Printf("  MIME检测失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ HTML文件MIME: %s\n", mimeType)
	}

	// 清理
	fs.Remove("test_demo.html")
}

// 测试文件复制
func testFileCopy(fs storage.Storage) {
	// 创建源文件
	fs.WriteFile("source_demo.txt", []byte("Source Content"), 0644)

	// 复制文件
	err := fs.Copy("source_demo.txt", "dest_demo.txt")
	if err != nil {
		fmt.Printf("  文件复制失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ 文件复制成功: source_demo.txt -> dest_demo.txt")

	// 验证复制结果
	data, _ := fs.ReadFile("dest_demo.txt")
	fmt.Printf("  ✓ 复制内容验证: %s\n", string(data))

	// 清理
	fs.Remove("source_demo.txt")
	fs.Remove("dest_demo.txt")
}

// 测试Excel操作
func testExcelOps(fs storage.Storage) {
	// 创建Excel文件
	file := fs.CreateExcel()

	// 设置单元格值
	file.SetCellValue("Sheet1", "A1", "姓名")
	file.SetCellValue("Sheet1", "B1", "年龄")
	file.SetCellValue("Sheet1", "A2", "张三")
	file.SetCellValue("Sheet1", "B2", 25)
	file.SetCellValue("Sheet1", "A3", "李四")
	file.SetCellValue("Sheet1", "B3", 30)

	// 保存Excel
	err := fs.SaveExcel(file, "demo.xlsx")
	if err != nil {
		fmt.Printf("  保存Excel失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ Excel文件创建成功: demo.xlsx")

	// 读取Excel
	rows, err := fs.ReadExcelSheet("demo.xlsx", "Sheet1")
	if err != nil {
		fmt.Printf("  读取Excel失败: %v\n", err)
		return
	}
	fmt.Printf("  ✓ Excel数据读取成功: %d 行\n", len(rows))
	for i, row := range rows {
		if i < 3 { // 只显示前3行
			fmt.Printf("    行 %d: %v\n", i+1, row)
		}
	}

	// 清理
	fs.Remove("demo.xlsx")
}

// 测试图片操作
func testImageOps(fs storage.Storage) {
	// 创建一个简单的测试图片 (纯色)
	img := imaging.New(200, 100, color.White)

	// 保存图片
	err := fs.SaveImage(img, "test_demo.png", imaging.PNG)
	if err != nil {
		fmt.Printf("  保存图片失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ 图片创建成功: test_demo.png")

	// 调整图片大小
	err = fs.ResizeImage("test_demo.png", "test_demo_resized.png", 100, 50, imaging.PNG)
	if err != nil {
		fmt.Printf("  调整图片大小失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ 图片调整大小成功: 200x100 -> 100x50")

	// 裁剪图片
	rect := image.Rect(10, 10, 90, 40)
	err = fs.CropImage("test_demo.png", "test_demo_cropped.png", rect, imaging.PNG)
	if err != nil {
		fmt.Printf("  裁剪图片失败: %v\n", err)
		return
	}
	fmt.Println("  ✓ 图片裁剪成功")

	// 清理
	fs.Remove("test_demo.png")
	fs.Remove("test_demo_resized.png")
	fs.Remove("test_demo_cropped.png")
}

// 测试文件监听
func testFileWatch(fs storage.Storage) {
	fmt.Println("  ℹ 文件监听功能需要实际运行环境测试")
	fmt.Println("  示例代码:")
	fmt.Println("    err := fs.Watch(\"./watch_dir\", func(event storage.WatchEvent) {")
	fmt.Println("        fmt.Printf(\"[%s] %s: %s\\n\", event.Time, event.Op, event.Path)")
	fmt.Println("    })")
	fmt.Println("  ✓ 文件监听接口已就绪")
}
