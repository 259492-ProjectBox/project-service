package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type CalendarHandler interface {
	CreateCalendar(c *gin.Context)
}

type calendarHandler struct {
	calendarService services.CalendarService
}

func NewCalendarHandler(calendarService services.CalendarService) CalendarHandler {
	return &calendarHandler{
		calendarService: calendarService,
	}
}

// @Summary Create a new Calendar
// @Description Creates a new calendarof that major
// @Tags Calendar
// @Accept  json
// @Produce  json
// @Param project body dtos.CreateCalendarRequest true "Calendar Data"
// @Success 201 {object} dtos.CalendarResponse "Successfully created calendar"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /calendar [post]
func (h *calendarHandler) CreateCalendar(c *gin.Context) {
	req := &dtos.CreateCalendarRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := h.calendarService.CreateCalendar(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// @Summary Update an existing event
// @Description Updates an event by its ID with the provided data
// @Tags Calendar
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param event body models.Event true "Updated Event Data"
// @Success 200 {object} models.Event "Successfully updated event"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /events/{id} [put]
// func (h *calendarHandler) UpdateEvent(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
// 		return
// 	}
// 	event := &models.Event{}
// 	if err := c.ShouldBindJSON(&event); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	_, err = h.calendarService.UpdateEvent(c, id, event)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, *event)
// }

// GetEventByID retrieves an event by its ID
// @Summary Get an event by ID
// @Description Fetches an event by its unique ID
// @Tags Calendar
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} models.Event "Successfully retrieved event"
// @Failure 400 {object} map[string]interface{} "Invalid event ID"
// @Failure 404 {object} map[string]interface{} "Event not found"
// @Router /events/{id} [get]
// func (h *calendarHandler) GetEventById(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
// 		return
// 	}

// 	event, err := h.calendarService.GetEventById(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, *event)
// }

// GetEventsByUserID retrieves events by User ID
// @Summary Get events by User ID
// @Description Fetches events by user ID
// @Tags Calendar
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} []models.Event "Successfully retrieved events"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "Events not found"
// @Router /events/user/{user_id} [get]
// func (h *calendarHandler) GetEventsByUserId(c *gin.Context) {
// 	userId := c.Param("user_id")

// 	events, err := h.calendarService.GetEventsByUserId(c, userId)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Events not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, events)
// }

// DeleteEvent deletes an event by its ID
// @Summary Delete an event by ID
// @Description Deletes the specified event using its ID
// @Tags Calendar
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Event deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /events/{id} [delete]
// func (h *calendarHandler) DeleteEvent(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
// 		return
// 	}

// 	if err := h.calendarService.DeleteEvent(c, id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
// }
