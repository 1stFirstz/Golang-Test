package services

import (
	"bewell_test/constants"
	"bewell_test/models"
	"regexp"
	"strconv"
	"strings"
)

type parsedProduct struct {
	ProductId  string
	MaterialId string
	ModelId    string
	Texture    string
	Qty        int
}

func splitProductId(cleanId string) (string, string) {
	parts := strings.Split(cleanId, constants.ProductIdSeparator)
	if len(parts) < constants.MinProductParts {
		return cleanId, ""
	}
	materialId := parts[0] + constants.ProductIdSeparator + parts[1]
	modelId := strings.Join(parts[2:], constants.ProductIdSeparator)
	return materialId, modelId
}

func CleanOrders(input []models.InputOrder) []models.CleanedOrder {
	var output []models.CleanedOrder
	currentNo := 1

	for _, order := range input {
		if strings.Contains(order.PlatformProductId, constants.BundleSeparator) {
			cleaned, nextNo := handleBundleOrder(order, currentNo)
			output = append(output, cleaned...)
			currentNo = nextNo
			continue
		}
		cleanId := cleanProductId(order.PlatformProductId)

		parts := strings.Split(cleanId, constants.ProductIdSeparator)
		if len(parts) < constants.MinProductParts {
			continue
		}
		materialId := parts[0] + constants.ProductIdSeparator + parts[1]
		modelId := strings.Join(parts[2:], constants.ProductIdSeparator)

		multiplier := extractQtyMultiplier(order.PlatformProductId)
		realQty := multiplier * order.Qty
		unitPrice := order.TotalPrice / float64(realQty)

		output = append(output, models.CleanedOrder{
			No:         currentNo,
			ProductId:  cleanId,
			MaterialId: materialId,
			ModelId:    modelId,
			Qty:        realQty,
			UnitPrice:  unitPrice,
			TotalPrice: order.TotalPrice,
		})
		currentNo++

		output = append(output, models.CleanedOrder{
			No:         currentNo,
			ProductId:  constants.WipingClothProductId,
			Qty:        realQty,
			UnitPrice:  constants.DefaultUnitPrice,
			TotalPrice: constants.DefaultTotalPrice,
		})
		currentNo++

		texture := parts[constants.MaterialIdPartIndex]
		cleaner := texture + constants.CleanerSuffix
		output = append(output, models.CleanedOrder{
			No:         currentNo,
			ProductId:  cleaner,
			Qty:        realQty,
			UnitPrice:  constants.DefaultUnitPrice,
			TotalPrice: constants.DefaultTotalPrice,
		})
		currentNo++
	}

	return consolidateOrders(output)
}

func cleanProductId(id string) string {
	re := regexp.MustCompile(constants.ProductIdPattern)
	matches := re.FindStringSubmatch(id)
	if len(matches) >= constants.MinRegexMatches {
		return matches[1]
	}
	return strings.TrimSpace(id)
}

func extractQtyMultiplier(id string) int {
	re := regexp.MustCompile(constants.QtyMultiplierPattern)
	matches := re.FindStringSubmatch(id)
	if len(matches) == constants.QtyRegexMatches {
		qty, err := strconv.Atoi(matches[1])
		if err == nil {
			return qty
		}
	}
	return constants.DefaultQtyMultiplier
}

func handleBundleOrder(order models.InputOrder, currentNo int) ([]models.CleanedOrder, int) {
	parts := strings.Split(order.PlatformProductId, constants.BundleSeparator)
	var result []models.CleanedOrder
	var items []parsedProduct
	textureQty := map[string]int{}
	textureOrder := []string{}
	totalQty := 0

	for _, part := range parts {
		cleanId := cleanProductId(part)
		matId, modelId := splitProductId(cleanId)
		texture := strings.Split(matId, constants.ProductIdSeparator)[constants.MaterialIdPartIndex]
		multiplier := extractQtyMultiplier(part)
		qty := multiplier * order.Qty
		totalQty += qty

		items = append(items, parsedProduct{
			ProductId:  cleanId,
			MaterialId: matId,
			ModelId:    modelId,
			Texture:    texture,
			Qty:        qty,
		})

		if _, exists := textureQty[texture]; !exists {
			textureOrder = append(textureOrder, texture)
		}
		textureQty[texture] += qty
	}

	unitPrice := constants.DefaultUnitPrice
	if totalQty > 0 {
		unitPrice = order.TotalPrice / float64(totalQty)
	}

	for _, item := range items {
		result = append(result, models.CleanedOrder{
			No:         currentNo,
			ProductId:  item.ProductId,
			MaterialId: item.MaterialId,
			ModelId:    item.ModelId,
			Qty:        item.Qty,
			UnitPrice:  unitPrice,
			TotalPrice: unitPrice * float64(item.Qty),
		})
		currentNo++
	}

	result = append(result, models.CleanedOrder{
		No:         currentNo,
		ProductId:  constants.WipingClothProductId,
		Qty:        totalQty,
		UnitPrice:  constants.DefaultUnitPrice,
		TotalPrice: constants.DefaultTotalPrice,
	})
	currentNo++

	for _, texture := range textureOrder {
		qty := textureQty[texture]
		result = append(result, models.CleanedOrder{
			No:         currentNo,
			ProductId:  texture + constants.CleanerSuffix,
			Qty:        qty,
			UnitPrice:  constants.DefaultUnitPrice,
			TotalPrice: constants.DefaultTotalPrice,
		})
		currentNo++
	}

	return result, currentNo
}

func consolidateOrders(orders []models.CleanedOrder) []models.CleanedOrder {
	consolidated := make(map[string]*models.CleanedOrder)
	var result []models.CleanedOrder
	var mainOrders []models.CleanedOrder

	for _, order := range orders {
		if order.ProductId == constants.WipingClothProductId ||
			strings.HasSuffix(order.ProductId, constants.CleanerSuffix) {
			if existing, exists := consolidated[order.ProductId]; exists {
				existing.Qty += order.Qty
			} else {
				consolidated[order.ProductId] = &models.CleanedOrder{
					ProductId:  order.ProductId,
					MaterialId: order.MaterialId,
					ModelId:    order.ModelId,
					Qty:        order.Qty,
					UnitPrice:  order.UnitPrice,
					TotalPrice: order.TotalPrice,
				}
			}
		} else {
			mainOrders = append(mainOrders, order)
		}
	}

	result = append(result, mainOrders...)

	currentNo := len(mainOrders) + 1

	if wipingCloth, exists := consolidated[constants.WipingClothProductId]; exists {
		wipingCloth.No = currentNo
		result = append(result, *wipingCloth)
		currentNo++
	}

	for _, texture := range constants.TextureOrder {
		cleanerName := texture + constants.CleanerSuffix
		if cleaner, exists := consolidated[cleanerName]; exists {
			cleaner.No = currentNo
			result = append(result, *cleaner)
			currentNo++
		}
	}

	for i := range result {
		result[i].No = i + 1
	}

	return result
}
