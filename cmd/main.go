package main

import (
	"bewell_test/models"
	"bewell_test/services"
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("🚀 Order Cleaning System Demo")
	fmt.Println("====================================================")

	// Test Case 1: Basic Order
	fmt.Println("\n📦 Test Case 1: Basic Order")
	basicOrder := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         50,
			TotalPrice:        100,
		},
	}
	processAndDisplay("Basic Order", basicOrder)

	// Test Case 2: Dirty Product ID
	fmt.Println("\n🧹 Test Case 2: Dirty Product ID")
	dirtyOrder := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX",
			Qty:               1,
			UnitPrice:         80,
			TotalPrice:        80,
		},
	}
	processAndDisplay("Dirty Product ID", dirtyOrder)

	// Test Case 3: Quantity Multiplier
	fmt.Println("\n✖️ Test Case 3: Quantity Multiplier")
	multiplierOrder := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-MATTE-IPHONE16PROMAX*3",
			Qty:               1,
			UnitPrice:         90,
			TotalPrice:        90,
		},
	}
	processAndDisplay("Quantity Multiplier", multiplierOrder)

	// Test Case 4: Bundle Order
	fmt.Println("\n📦📦 Test Case 4: Bundle Order")
	bundleOrder := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-OPPOA3/FG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
			Qty:               1,
			UnitPrice:         120,
			TotalPrice:        120,
		},
	}
	processAndDisplay("Bundle Order", bundleOrder)

	// Test Case 5: Complex Mixed Order
	fmt.Println("\n🎯 Test Case 5: Mixed Orders")
	mixedOrders := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE15PRO*2",
			Qty:               1,
			UnitPrice:         60,
			TotalPrice:        60,
		},
		{
			No:                2,
			PlatformProductId: "noise&FG0A-MATTE-GALAXYS24",
			Qty:               3,
			UnitPrice:         45,
			TotalPrice:        135,
		},
	}
	processAndDisplay("Mixed Orders", mixedOrders)

	fmt.Println("\n🎉 Demo completed successfully!")
}

func processAndDisplay(title string, orders []models.InputOrder) {
	fmt.Printf("\n--- %s ---\n", title)

	// Display input
	fmt.Println("📥 Input Orders:")
	for _, order := range orders {
		fmt.Printf("  • Order #%d: %s (Qty: %d, Total: %.2f)\n",
			order.No, order.PlatformProductId, order.Qty, order.TotalPrice)
	}

	// Process orders
	cleaned := services.CleanOrders(orders)

	// Display output
	fmt.Printf("\n📤 Cleaned Orders (%d items):\n", len(cleaned))
	for _, order := range cleaned {
		status := "🔧"
		if order.ProductId == "WIPING-CLOTH" {
			status = "🧽"
		} else if len(order.ProductId) > 8 && order.ProductId[len(order.ProductId)-8:] == "CLEANNER" {
			status = "🧴"
		}

		fmt.Printf("  %s Order #%d: %s\n", status, order.No, order.ProductId)
		if order.MaterialId != "" {
			fmt.Printf("      Material: %s, Model: %s\n", order.MaterialId, order.ModelId)
		}
		fmt.Printf("      Qty: %d, Unit: %.2f, Total: %.2f\n",
			order.Qty, order.UnitPrice, order.TotalPrice)
		fmt.Println()
	}

	fmt.Println("📋 JSON Output:")
	jsonData, _ := json.MarshalIndent(cleaned, "", "  ")
	fmt.Println(string(jsonData))
}
