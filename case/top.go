package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {

	const dataSize = 1000000000 // 10亿数据

	fmt.Printf("开始生成 %d 条测试数据...\n", dataSize)
	startTime := time.Now()
	data := generateData(dataSize)
	fmt.Printf("数据生成完成，耗时: %v\n", time.Since(startTime))

	fmt.Printf("\n开始并行查找最大值...\n")
	startTime = time.Now()
	maxParallel := findMaxParallel(data)
	parallelTime := time.Since(startTime)
	fmt.Printf("并行查找结果: %d, 耗时: %v\n", maxParallel, parallelTime)

	fmt.Printf("\n开始顺序查找最大值...\n")
	startTime = time.Now()
	maxSequential := findMaxSequential(data)
	sequentialTime := time.Since(startTime)
	fmt.Printf("顺序查找结果: %d, 耗时: %v\n", maxSequential, sequentialTime)

	fmt.Printf("\n性能对比:\n")
	fmt.Printf("并行查找耗时: %v\n", parallelTime)
	fmt.Printf("顺序查找耗时: %v\n", sequentialTime)
	fmt.Printf("加速比: %.2fx\n", float64(sequentialTime.Nanoseconds())/float64(parallelTime.Nanoseconds()))

	// 验证结果一致性
	if maxParallel == maxSequential {
		fmt.Printf("✓ 两种方法结果一致\n")
	} else {
		fmt.Printf("✗ 结果不一致: 并行=%d, 顺序=%d\n", maxParallel, maxSequential)
	}

}

// generateData 生成测试数据
func generateData(size int) []int {
	data := make([]int, size)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	for i := range data {
		data[i] = rng.Intn(1000_000_000)
	}
	// 设置一个明确的最大值
	data[size/2] = 2000000000
	return data
}

// findMaxInRange 在指定范围内查找最大值
func findMaxInRange(data []int, start, end int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	max := data[start]
	for i := start + 1; i < end; i++ {
		if max < data[i] {
			max = data[i]
		}
	}
	result <- max
}

// findMaxParallel 并行查找最大值
func findMaxParallel(data []int) int {
	n := len(data)
	if n == 0 {
		return 0
	}

	numCPU := runtime.NumCPU()
	chunkSize := n / numCPU
	resultChan := make(chan int, numCPU)
	var wg sync.WaitGroup

	// [start，end) 中查找
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numCPU-1 {
			end = n
		}
		wg.Add(1)
		go findMaxInRange(data, start, end, resultChan, &wg)
	}
	wg.Wait()
	close(resultChan)

	// 收集各分区的最大值，找出全局最大值
	finalMax := <-resultChan
	for max := range resultChan {
		if max > finalMax {
			finalMax = max
		}
	}
	return finalMax
}

// findMaxSequential 顺序查找最大值（用于对比）
func findMaxSequential(data []int) int {
	if len(data) == 0 {
		return 0
	}

	max := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
		}
	}
	return max
}
