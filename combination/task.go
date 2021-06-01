package combination

import (
	"fmt"
	"net/http"
	"time"
)

//CheckoutPageComponent 订单结算页面组件
type CheckoutPageComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *CheckoutPageComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 订单结算页面组件...")
	return
}

//AddressComponent 地址组件
type AddressComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *AddressComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 地址组件...")
	fmt.Println(runFuncName(), " 获取地址信息 ing...")

	//模拟远程调用地址服务
	http.Get("http://example.com/")
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取地址信息 done...")
	return
}

//PayMethodComponent 支付方式组件
type PayMethodComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *PayMethodComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 支付方式组件...")
	fmt.Println(runFuncName(), " 获取支付方式 ing...")

	//模拟远程调用地址服务 略
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取支付方式 done...")
	return
}

//StoreComponent 店铺组件
type StoreComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *StoreComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 店铺组件...")
	return
}

//SkuComponent 商品组件
type SkuComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *SkuComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 商品组件...")
	return
}

//PromotionComponent 优惠信息组件
type PromotionComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *PromotionComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 优惠信息组件...")
	return
}

//ExpressComponent 物流组件
type ExpressComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *ExpressComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 物流组件...")
	return
}

//AftersaleComponent 售后组件
type AftersaleComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *AftersaleComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 售后组件...")
	return
}

//InvoiceComponent 发票组件
type InvoiceComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *InvoiceComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 发票组件...")
	fmt.Println(runFuncName(), " 获取发票信息 ing...")
	//模拟远程调用地址服务 略
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取发票信息 done...")
	return
}

//CouponComponent 优惠券组件
type CouponComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *CouponComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 优惠券组件...")
	fmt.Println(runFuncName(), " 获取优惠券信息 ing...")
	//模拟远程调用地址服务
	http.Get("http://example.com/")
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取最优优惠券 done...")
	return
}

//GiftCardComponent 礼品卡组件
type GiftCardComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际业务逻辑代码填充处
func (bc *GiftCardComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 礼品卡组件...")
	fmt.Println(runFuncName(), " 获取礼品卡信息 ing...")
	//模拟远程调用地址服务
	http.Get("http://example.com/")
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取礼品卡信息 done...")
	return
}

// OrderComponent 订单金额详细信息组件
type OrderComponent struct {
	// 合成复用基础组件
	BaseComponent
}

//BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *OrderComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), " 订单金额详细信息组件...")
	return
}

// Demo 示例
func Demo() {
	//初始化订单结算洁面这个大组件
	checkoutPage := &CheckoutPageComponent{}

	// 挂载子组件
	storeComponent := &StoreComponent{}
	skuComponent := &SkuComponent{}
	skuComponent.Mount(&PromotionComponent{}, &AftersaleComponent{})
	storeComponent.Mount(skuComponent, &ExpressComponent{})

	// ---挂载组件---

	//普通组件
	checkoutPage.Mount(storeComponent, &OrderComponent{})

	//并发组件
	checkoutPage.MountConcurrency(
		&AddressComponent{},
		&PayMethodComponent{},
		&InvoiceComponent{},
		&CouponComponent{},
		&GiftCardComponent{},
	)
	//初始化业务上下文 并设置超时时间
	ctx := GetContext(5 * time.Second)
	defer ctx.CancelFunc()
	//开始构建页面组件数据
	checkoutPage.ChildsDo(ctx)
}
