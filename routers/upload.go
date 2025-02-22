package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupUploadRouter(r *gin.RouterGroup, handler handlers.UploadHandler) {
	uploadRouteV1 := r.Group("/v1/uploads")
	{
		uploadRouteV1.POST("/program/:program_id/student-enrollment", handler.UploadStudentEnrollmentFile)
		uploadRouteV1.POST("/program/:program_id/create-project", handler.UploadCreateProjectFile)
		uploadRouteV1.POST("/program/:program_id/create-staff", handler.UploadCreateStaffFile)
	}
}
