
/*
    - /api/randomFilteredImage/{image_link:path}/{kernel_size}/{low}/{high}/{kernel_type}
    - /api/invertedImage/{image_link:path}
    - /api/saturatedImage/{image_link:path}/{saturation}
    - /api/edgeImage/{image_link:path}/{lower}/{higher}
    - /api/dilatedImage/{image_link:path}/{box_size}/{iterations}
    - /api/erodedImage/{image_link:path}/{box_size}/{iterations}
    - /api/textImage/{image_link:path}/{text}/{font_scale}/{x}/{y}
    - /api/reducedImage/{image_link:path}/{quality}
    - /api/shuffledImage/{image_link:path}/{partitions}
*/

package imagemanipulation




import (
	"github.com/labstack/echo/v4"
)


func init_routing(e  *echo.Echo){}

func main(){
e := echo.New()

init_routing(e)
}
