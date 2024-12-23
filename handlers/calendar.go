package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type CalendarHandler interface {
	CreateCalendarHandler(c *gin.Context)
	GetCalendarByMajorIDHandler(c *gin.Context)
	UpdateCalendarHandler(c *gin.Context)
	DeleteCalendarHandler(c *gin.Context)
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
func (h *calendarHandler) CreateCalendarHandler(c *gin.Context) {
	req := &dtos.CreateCalendarRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := h.calendarService.CreateCalendarService(c, req)
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

// @Summary Get calendar by major ID
// @Description Fetches all calendar events for a given major
// @Tags Calendar
// @Produce  json
// @Param major_id path int true "Major ID"
// @Success 200 {object} []dtos.CalendarResponse "Successfully retrieved calendar events"
// @Failure 400 {object} map[string]interface{} "Invalid major ID"
// @Failure 404 {object} map[string]interface{} "Calendar events not found"
// @Router /calendar/GetByMajorID/{major_id} [get]
func (h *calendarHandler) GetCalendarByMajorIDHandler(c *gin.Context) {
	majorID, err := strconv.Atoi(c.Param("major_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid major ID"})
		return
	}
	events, err := h.calendarService.GetCalendarByMajorIDService(c, majorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar events not found"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// @Summary Update an existing event
// @Description Updates an event by its ID with the provided data
// @Tags Calendar
// @Accept  json
// @Produce  json
// @Param event body dtos.UpdateCalendarRequest true "Updated Event Data"
// @Success 200 {object} dtos.CalendarResponse "Successfully updated event"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /calendar [put]
func (h *calendarHandler) UpdateCalendarHandler(c *gin.Context) {
	calendar := &dtos.UpdateCalendarRequest{}
	if err := c.ShouldBindJSON(&calendar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCalendar, err := h.calendarService.UpdateCalendarService(c, calendar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCalendar)
}

// DeleteEvent deletes an event by its ID
// @Summary Delete an event by ID
// @Description Deletes the specified event using its ID
// @Tags Calendar
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Event deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /calendar/{id} [delete]
func (h *calendarHandler) DeleteCalendarHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := h.calendarService.DeleteCalendarService(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
