package handlers

import (
	"strconv"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/services"
	"url-shortener-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	urlService *services.URLService
}

func NewURLHandler(urlService *services.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) CreateURL(c *fiber.Ctx) error {
	var req models.CreateURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}
	
	var userID *uint
	if uid := c.Locals("user_id"); uid != nil {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}
	
	url, err := h.urlService.CreateURL(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "url_creation_failed",
			Message: err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Data:    url,
		Message: "URL created successfully",
	})
}

func (h *URLHandler) RedirectURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	
	url, err := h.urlService.GetURLByShortCode(shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "url_not_found",
			Message: err.Error(),
		})
	}
	
	analytics := &models.Analytics{
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
		Referrer:  c.Get("Referer"),
	}
	
	device, os, browser := utils.ParseUserAgent(analytics.UserAgent)
	analytics.Device = device
	analytics.OS = os
	analytics.Browser = browser
	
	go h.urlService.RecordClick(url.ID, analytics)
	
	return c.Redirect(url.OriginalURL, fiber.StatusFound)
}

func (h *URLHandler) GetURLInfo(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	
	url, err := h.urlService.GetURLByShortCode(shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "url_not_found",
			Message: err.Error(),
		})
	}
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data:    url,
	})
}

func (h *URLHandler) GetUserURLs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	
	if limit > 100 {
		limit = 100
	}
	
	urls, total, err := h.urlService.GetUserURLs(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"urls":  urls,
			"total": total,
			"limit": limit,
			"offset": offset,
		},
	})
}

func (h *URLHandler) UpdateURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	urlID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_url_id",
			Message: "Invalid URL ID",
		})
	}
	
	var req models.CreateURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}
	
	url, err := h.urlService.UpdateURL(uint(urlID), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
	}
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data:    url,
		Message: "URL updated successfully",
	})
}

func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	urlID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_url_id",
			Message: "Invalid URL ID",
		})
	}
	
	if err := h.urlService.DeleteURL(uint(urlID), userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
	}
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Message: "URL deleted successfully",
	})
}

func (h *URLHandler) GetURLAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	
	urlID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_url_id",
			Message: "Invalid URL ID",
		})
	}
	
	analytics, err := h.urlService.GetURLAnalytics(uint(urlID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}
	
	stats, err := h.urlService.GetURLStats(uint(urlID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"analytics": analytics,
			"stats":     stats,
		},
	})
}