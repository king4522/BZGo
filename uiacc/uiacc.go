package uiacc

import "C"
import (
	"strconv"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/utils"

	"github.com/Dasongzi1366/AutoGo/java"
)

type Uiacc struct {
	index    int
	selector string
}

type UiObject struct {
	objStr string
}

type Rect struct {
	Left    int
	Right   int
	Top     int
	Bottom  int
	CenterX int
	CenterY int
	Width   int
	Height  int
}

var index int

func init() {
	utils.Shell("pkill -f assistdir")
	java.CallJavaMethod("acc", "newAccessibility")
}

// New 创建一个新的 Accessibility 对象
func New() *Uiacc {
	index++
	return &Uiacc{index: index}
}

// Text 设置选择器的 text 属性
func (a *Uiacc) Text(value string) *Uiacc {
	a.selector = a.selector + "text@@" + value + "&&"
	return a
}

// TextContains 设置选择器的 textContains 属性，用于匹配包含指定文本的控件
func (a *Uiacc) TextContains(value string) *Uiacc {
	a.selector = a.selector + "textContains@@" + value + "&&"
	return a
}

// TextStartsWith 设置选择器的 textStartsWith 属性，用于匹配以指定文本开头的控件
func (a *Uiacc) TextStartsWith(value string) *Uiacc {
	a.selector = a.selector + "textStartsWith@@" + value + "&&"
	return a
}

// TextEndsWith 设置选择器的 textEndsWith 属性，用于匹配以指定文本结尾的控件
func (a *Uiacc) TextEndsWith(value string) *Uiacc {
	a.selector = a.selector + "textEndsWith@@" + value + "&&"
	return a
}

// TextMatches 设置选择器的 textMatches 属性，用于匹配符合指定正则表达式的控件
func (a *Uiacc) TextMatches(value string) *Uiacc {
	a.selector = a.selector + "textMatches@@" + value + "&&"
	return a
}

// Desc 设置选择器的 desc 属性，用于匹配描述等于指定文本的控件
func (a *Uiacc) Desc(value string) *Uiacc {
	a.selector = a.selector + "desc@@" + value + "&&"
	return a
}

// DescContains 设置选择器的 descContains 属性，用于匹配描述包含指定文本的控件
func (a *Uiacc) DescContains(value string) *Uiacc {
	a.selector = a.selector + "descContains@@" + value + "&&"
	return a
}

// DescStartsWith 设置选择器的 descStartsWith 属性，用于匹配描述以指定文本开头的控件
func (a *Uiacc) DescStartsWith(value string) *Uiacc {
	a.selector = a.selector + "descStartsWith@@" + value + "&&"
	return a
}

// DescEndsWith 设置选择器的 descEndsWith 属性，用于匹配描述以指定文本结尾的控件
func (a *Uiacc) DescEndsWith(value string) *Uiacc {
	a.selector = a.selector + "descEndsWith@@" + value + "&&"
	return a
}

// DescMatches 设置选择器的 descMatches 属性，用于匹配描述符合指定正则表达式的控件
func (a *Uiacc) DescMatches(value string) *Uiacc {
	a.selector = a.selector + "descMatches@@" + value + "&&"
	return a
}

// Id 设置选择器的 id 属性，用于匹配ID等于指定值的控件
func (a *Uiacc) Id(value string) *Uiacc {
	a.selector = a.selector + "id@@" + value + "&&"
	return a
}

// IdContains 设置选择器的 idContains 属性，用于匹配ID包含指定值的控件
func (a *Uiacc) IdContains(value string) *Uiacc {
	a.selector = a.selector + "idContains@@" + value + "&&"
	return a
}

// IdStartsWith 设置选择器的 idStartsWith 属性，用于匹配ID以指定值开头的控件
func (a *Uiacc) IdStartsWith(value string) *Uiacc {
	a.selector = a.selector + "idStartsWith@@" + value + "&&"
	return a
}

// IdEndsWith 设置选择器的 idEndsWith 属性，用于匹配ID以指定值结尾的控件
func (a *Uiacc) IdEndsWith(value string) *Uiacc {
	a.selector = a.selector + "idEndsWith@@" + value + "&&"
	return a
}

// IdMatches 设置选择器的 idMatches 属性，用于匹配ID符合指定正则表达式的控件
func (a *Uiacc) IdMatches(value string) *Uiacc {
	a.selector = a.selector + "idMatches@@" + value + "&&"
	return a
}

// ClassName 设置选择器的 className 属性，用于匹配类名等于指定值的控件
func (a *Uiacc) ClassName(value string) *Uiacc {
	a.selector = a.selector + "className@@" + value + "&&"
	return a
}

// ClassNameContains 设置选择器的 classNameContains 属性，用于匹配类名包含指定值的控件
func (a *Uiacc) ClassNameContains(value string) *Uiacc {
	a.selector = a.selector + "classNameContains@@" + value + "&&"
	return a
}

// ClassNameStartsWith 设置选择器的 classNameStartsWith 属性，用于匹配类名以指定值开头的控件
func (a *Uiacc) ClassNameStartsWith(value string) *Uiacc {
	a.selector = a.selector + "classNameStartsWith@@" + value + "&&"
	return a
}

// ClassNameEndsWith 设置选择器的 classNameEndsWith 属性，用于匹配类名以指定值结尾的控件
func (a *Uiacc) ClassNameEndsWith(value string) *Uiacc {
	a.selector = a.selector + "classNameEndsWith@@" + value + "&&"
	return a
}

// ClassNameMatches 设置选择器的 classNameMatches 属性，用于匹配类名符合指定正则表达式的控件
func (a *Uiacc) ClassNameMatches(value string) *Uiacc {
	a.selector = a.selector + "classNameMatches@@" + value + "&&"
	return a
}

// PackageName 设置选择器的 packageName 属性，用于匹配包名等于指定值的控件
func (a *Uiacc) PackageName(value string) *Uiacc {
	a.selector = a.selector + "packageName@@" + value + "&&"
	return a
}

// PackageNameContains 设置选择器的 packageNameContains 属性，用于匹配包名包含指定值的控件
func (a *Uiacc) PackageNameContains(value string) *Uiacc {
	a.selector = a.selector + "packageNameContains@@" + value + "&&"
	return a
}

// PackageNameStartsWith 设置选择器的 packageNameStartsWith 属性，用于匹配包名以指定值开头的控件
func (a *Uiacc) PackageNameStartsWith(value string) *Uiacc {
	a.selector = a.selector + "packageNameStartsWith@@" + value + "&&"
	return a
}

// PackageNameEndsWith 设置选择器的 packageNameEndsWith 属性，用于匹配包名以指定值结尾的控件
func (a *Uiacc) PackageNameEndsWith(value string) *Uiacc {
	a.selector = a.selector + "packageNameEndsWith@@" + value + "&&"
	return a
}

// PackageNameMatches 设置选择器的 packageNameMatches 属性，用于匹配包名符合指定正则表达式的控件
func (a *Uiacc) PackageNameMatches(value string) *Uiacc {
	a.selector = a.selector + "packageNameMatches@@" + value + "&&"
	return a
}

// Bounds 设置选择器的 bounds 属性，用于匹配控件在屏幕上的范围
func (a *Uiacc) Bounds(left, top, right, bottom int) *Uiacc {
	a.selector = a.selector + "bounds@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"
	return a
}

// BoundsInside 设置选择器的 boundsInside 属性，用于匹配控件在屏幕内的范围
func (a *Uiacc) BoundsInside(left, top, right, bottom int) *Uiacc {
	a.selector = a.selector + "boundsInside@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"
	return a
}

// BoundsContains 设置选择器的 boundsContains 属性，用于匹配控件包含在指定范围内
func (a *Uiacc) BoundsContains(left, top, right, bottom int) *Uiacc {
	a.selector = a.selector + "boundsContains@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"
	return a
}

// DrawingOrder 设置选择器的 drawingOrder 属性，用于匹配控件在父控件中的绘制顺序
func (a *Uiacc) DrawingOrder(value int) *Uiacc {
	a.selector = a.selector + "drawingOrder@@" + i2s(value) + "&&"
	return a
}

// Clickable 设置选择器的 clickable 属性，用于匹配控件是否可点击
func (a *Uiacc) Clickable(value bool) *Uiacc {
	a.selector = a.selector + "clickAble@@" + b2s(value) + "&&"
	return a
}

// LongClickable 设置选择器的 longClickable 属性，用于匹配控件是否可长按
func (a *Uiacc) LongClickable(value bool) *Uiacc {
	a.selector = a.selector + "longClickAble@@" + b2s(value) + "&&"
	return a
}

// Checkable 设置选择器的 checkable 属性，用于匹配控件是否可选中
func (a *Uiacc) Checkable(value bool) *Uiacc {
	a.selector = a.selector + "checkAble@@" + b2s(value) + "&&"
	return a
}

// Selected 设置选择器的 selected 属性，用于匹配控件是否被选中
func (a *Uiacc) Selected(value bool) *Uiacc {
	a.selector = a.selector + "selected@@" + b2s(value) + "&&"
	return a
}

// Enabled 设置选择器的 enabled 属性，用于匹配控件是否启用
func (a *Uiacc) Enabled(value bool) *Uiacc {
	a.selector = a.selector + "enabled@@" + b2s(value) + "&&"
	return a
}

// Scrollable 设置选择器的 scrollable 属性，用于匹配控件是否可滚动
func (a *Uiacc) Scrollable(value bool) *Uiacc {
	a.selector = a.selector + "scrollAble@@" + b2s(value) + "&&"
	return a
}

// Editable 设置选择器的 editable 属性，用于匹配控件是否可编辑
func (a *Uiacc) Editable(value bool) *Uiacc {
	a.selector = a.selector + "editable@@" + b2s(value) + "&&"
	return a
}

// MultiLine 设置选择器的 multiLine 属性，用于匹配控件是否多行
func (a *Uiacc) MultiLine(value bool) *Uiacc {
	a.selector = a.selector + "multiLine@@" + b2s(value) + "&&"
	return a
}

// Checked 设置选择器的 checked 属性，用于匹配控件是否被勾选
func (a *Uiacc) Checked(value bool) *Uiacc {
	a.selector = a.selector + "checked@@" + b2s(value) + "&&"
	return a
}

// Focusable 设置选择器的 focusable 属性，用于匹配控件是否可聚焦
func (a *Uiacc) Focusable(value bool) *Uiacc {
	a.selector = a.selector + "focusable@@" + b2s(value) + "&&"
	return a
}

// Dismissable 设置选择器的 dismissable 属性，用于匹配控件是否可解散
func (a *Uiacc) Dismissable(value bool) *Uiacc {
	a.selector = a.selector + "dismissable@@" + b2s(value) + "&&"
	return a
}

// Focused 设置选择器的 UiaccFocused 属性，用于匹配控件是否是辅助功能焦点
func (a *Uiacc) Focused(value bool) *Uiacc {
	a.selector = a.selector + "Focused@@" + b2s(value) + "&&"
	return a
}

// ContextClickable 设置选择器的 contextClickable 属性，用于匹配控件是否是上下文点击
func (a *Uiacc) ContextClickable(value bool) *Uiacc {
	a.selector = a.selector + "contextClickable@@" + b2s(value) + "&&"
	return a
}

// Index 设置选择器的 index 属性，用于匹配控件在父控件中的索引
func (a *Uiacc) Index(value int) *Uiacc {
	a.selector = a.selector + "indexInParent@@" + i2s(value) + "&&"
	return a
}

// Click 点击屏幕上的文本
func (a *Uiacc) Click(text string) bool {
	obj := a.Text(text).FindOnce()
	if obj != nil {
		return obj.Click() || obj.GetParent().Click()
	} else {
		obj := a.Desc(text).FindOnce()
		if obj != nil {
			return obj.Click() || obj.GetParent().Click()
		}
	}
	return false
}

// WaitFor 等待控件出现并返回 UiObject 对象 超时单位为毫秒,写0代表无限等待,超时返回nil
func (a *Uiacc) WaitFor(timeout int) *UiObject {
	var str string
	startTime := time.Now()
	for {
		str = java.CallJavaMethod("acc", "findOnce|"+a.selector)
		if str != "" {
			break
		}
		if timeout > 0 && time.Since(startTime).Milliseconds() >= int64(timeout) {
			a.selector = ""
			return nil
		}
		sleep(100)
	}
	a.selector = ""
	return &UiObject{objStr: str}
}

// FindOnce 查找单个控件并返回 UiObject 对象
func (a *Uiacc) FindOnce() *UiObject {
	str := java.CallJavaMethod("acc", "findOnce|"+a.selector)
	a.selector = ""
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// Find 查找所有符合条件的控件并返回 UiObject 对象数组
func (a *Uiacc) Find() []*UiObject {
	str := java.CallJavaMethod("acc", "find|"+a.selector)
	a.selector = ""
	if str == "" {
		return nil
	}
	arr := strings.Split(str, "\n")
	var uiObjectArray []*UiObject
	for _, s := range arr {
		if s != "" {
			uiObjectArray = append(uiObjectArray, &UiObject{objStr: s})
		}
	}
	return uiObjectArray
}

// FindS 查找所有符合条件的控件并返回文本字符串
func (a *Uiacc) FindS() string {
	str := java.CallJavaMethod("acc", "findS|"+a.selector)
	a.selector = ""
	return str
}

// Close 关闭无障碍服务
func (a *Uiacc) Close() {
	java.CallJavaMethod("acc", "close")
}

// Click 点击该控件，并返回是否点击成功
func (u *UiObject) Click() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectClick|"+u.objStr))
}

// ClickCenter 使用坐标点击该控件的中点，相当于click(uiObj.bounds().centerX(), uiObject.bounds().centerY())
func (u *UiObject) ClickCenter() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectClickCenter|"+u.objStr))
}

// ClickLongClick 长按该控件，并返回是否点击成功
func (u *UiObject) ClickLongClick() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectLongClick|"+u.objStr))
}

// Copy 对输入框文本的选中内容进行复制，并返回是否操作成功
func (u *UiObject) Copy() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectCopy|"+u.objStr))
}

// Cut 对输入框文本的选中内容进行剪切，并返回是否操作成功
func (u *UiObject) Cut() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectCut|"+u.objStr))
}

// Paste 对输入框控件进行粘贴操作，把剪贴板内容粘贴到输入框中，并返回是否操作成功
func (u *UiObject) Paste() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectPaste|"+u.objStr))
}

// ScrollForward 对控件执行向前滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollForward() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectScrollForward|"+u.objStr))
}

// ScrollBackward 对控件执行向后滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollBackward() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectScrollBackward|"+u.objStr))
}

// Collapse 对控件执行折叠操作，并返回是否操作成功
func (u *UiObject) Collapse() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectCollapse|"+u.objStr))
}

// Expand 对控件执行展开操作，并返回是否操作成功
func (u *UiObject) Expand() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectExpand|"+u.objStr))
}

// Show 执行显示操作，并返回是否操作成功
func (u *UiObject) Show() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectShow|"+u.objStr))
}

// Select 对控件执行"选中"操作，并返回是否操作成功
func (u *UiObject) Select() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectSelect|"+u.objStr))
}

// ClearSelect 清除控件的选中状态，并返回是否操作成功
func (u *UiObject) ClearSelect() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectClearSelect|"+u.objStr))
}

// SetSelection 对输入框控件设置选中的文字内容，并返回是否操作成功
func (u *UiObject) SetSelection(start, end int) bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectSetSelection|"+u.objStr+"|"+i2s(start)+"|"+i2s(end)))
}

// SetText 设置输入框控件的文本内容，并返回是否设置成功
func (u *UiObject) SetText(str string) bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectSetText|"+u.objStr+"|"+str))
}

// GetClickable 获取控件的 clickable 属性
func (u *UiObject) GetClickable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetClickAble|"+u.objStr))
}

// GetLongClickable 获取控件的 longClickable 属性
func (u *UiObject) GetLongClickable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetLongClickAble|"+u.objStr))
}

// GetCheckable 获取控件的 checkable 属性
func (u *UiObject) GetCheckable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetCheckable|"+u.objStr))
}

// GetSelected 获取控件的 selected 属性
func (u *UiObject) GetSelected() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetSelected|"+u.objStr))
}

// GetEnabled 获取控件的 enabled 属性
func (u *UiObject) GetEnabled() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetEnabled|"+u.objStr))
}

// GetScrollable 获取控件的 scrollable 属性
func (u *UiObject) GetScrollable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetScrollAble|"+u.objStr))
}

// GetEditable 获取控件的 editable 属性
func (u *UiObject) GetEditable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetEditable|"+u.objStr))
}

// GetMultiLine 获取控件的 multiLine 属性
func (u *UiObject) GetMultiLine() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetMultiLine|"+u.objStr))
}

// GetChecked 获取控件的 checked 属性
func (u *UiObject) GetChecked() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetChecked|"+u.objStr))
}

// GetFocusable 获取控件的 focusable 属性
func (u *UiObject) GetFocusable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetFocusable|"+u.objStr))
}

// GetDismissable 获取控件的 dismissable 属性
func (u *UiObject) GetDismissable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetDismissable|"+u.objStr))
}

// GetContextClickable 获取控件的 contextClickable 属性
func (u *UiObject) GetContextClickable() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetContextClickable|"+u.objStr))
}

// GetUiaccFocused 获取控件的 UiaccFocused 属性
func (u *UiObject) GetUiaccFocused() bool {
	return s2b(java.CallJavaMethod("acc", "uiObjectGetUiaccFocused|"+u.objStr))
}

// GetChildCount 获取控件的子控件数目
func (u *UiObject) GetChildCount() int {
	return s2i(java.CallJavaMethod("acc", "uiObjectGetChildCount|"+u.objStr))
}

// GetDrawingOrder 获取控件在父控件中的绘制次序
func (u *UiObject) GetDrawingOrder() int {
	return s2i(java.CallJavaMethod("acc", "uiObjectGetDrawingOrder|"+u.objStr))
}

// GetIndex 获取控件在父控件中的索引
func (u *UiObject) GetIndex() int {
	return s2i(java.CallJavaMethod("acc", "uiObjectGetIndexInParent|"+u.objStr))
}

// GetBounds 获取控件在屏幕上的范围
func (u *UiObject) GetBounds() Rect {
	bounds := java.CallJavaMethod("acc", "uiObjectGetBounds|"+u.objStr)
	arr := strings.Split(bounds, ",")
	if len(arr) == 4 {
		return Rect{Left: s2i(arr[0]), Top: s2i(arr[1]), Right: s2i(arr[2]), Bottom: s2i(arr[3]), Width: s2i(arr[2]) - s2i(arr[0]), Height: s2i(arr[3]) - s2i(arr[1]), CenterX: s2i(arr[0]) + ((s2i(arr[2]) - s2i(arr[0])) / 2), CenterY: s2i(arr[1]) + ((s2i(arr[3]) - s2i(arr[1])) / 2)}
	}
	return Rect{}
}

// GetBoundsInParent 获取控件在父控件中的范围
func (u *UiObject) GetBoundsInParent() string {
	return java.CallJavaMethod("acc", "uiObjectGetBoundsInParent|"+u.objStr)
}

// GetId 获取控件的ID
func (u *UiObject) GetId() string {
	return java.CallJavaMethod("acc", "uiObjectGetId|"+u.objStr)
}

// GetText 获取控件的文本内容
func (u *UiObject) GetText() string {
	return java.CallJavaMethod("acc", "uiObjectGetText|"+u.objStr)
}

// GetDesc 获取控件的描述内容
func (u *UiObject) GetDesc() string {
	return java.CallJavaMethod("acc", "uiObjectGetDesc|"+u.objStr)
}

// GetPackageName 获取控件的包名
func (u *UiObject) GetPackageName() string {
	return java.CallJavaMethod("acc", "uiObjectGetPackageName|"+u.objStr)
}

// GetClassName 获取控件的类名
func (u *UiObject) GetClassName() string {
	return java.CallJavaMethod("acc", "uiObjectGetClassName|"+u.objStr)
}

// GetParent 获取控件的父控件
func (u *UiObject) GetParent() *UiObject {
	str := java.CallJavaMethod("acc", "uiObjectGetParent|"+u.objStr)
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// GetChild 获取控件的指定索引的子控件
func (u *UiObject) GetChild(index int) *UiObject {
	str := java.CallJavaMethod("acc", "uiObjectGetChild|"+u.objStr+"|"+i2s(index))
	if str == "" {
		return nil
	}
	return &UiObject{objStr: str}
}

// GetChildren 获取控件的所有子控件
func (u *UiObject) GetChildren() []*UiObject {
	str := java.CallJavaMethod("acc", "uiObjectGetChildren|"+u.objStr)
	if str == "" {
		return nil
	}
	arr := strings.Split(str, "\n")
	var uiObjectArray []*UiObject
	for _, s := range arr {
		if s != "" {
			uiObjectArray = append(uiObjectArray, &UiObject{objStr: s})
		}
	}
	return uiObjectArray
}

// ToString 节点对象转文本
func (u *UiObject) ToString() string {
	return java.CallJavaMethod("acc", "uiObjectToString|"+u.objStr)
}

func b2s(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func s2b(s string) bool {
	return s == "true"
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}
