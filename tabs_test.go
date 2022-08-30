//

package tabs

import (
	"fmt"
	"testing"
	"time"

	"github.com/dudnikv/pro"
	"github.com/dudnikv/ru"
)

func TestTabs(t *testing.T) {

	fmt.Println(pro.Nul, pro.Any, pro.All)

	var teas easVersionType
	fmt.Println(teas.Parse(".20.1.0.13623"))

	xx, _ := teas.Parse(".20.1.0.13623")

	fmt.Println(teas.Emit(xx))

	fmt.Println(NewEnumValueType("bool", "true", "false"))
	fmt.Println(registeredValueTypes)
	fmt.Println(ru.DateTime(time.Now()))
	//	PingList("r78-188230-n", "R78Dudnik", "R78Dudnik1", "10.47.25.62", "10.46.12.104")
}
