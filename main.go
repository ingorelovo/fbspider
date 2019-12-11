// https://www.w3schools.com/cssref/css_selectors.asp
package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	// These paths will be different on your system.
	port          = 8080
	timeDelta     = time.Second * 5
	registerIndex = "https://m.myweimai.com/wx/dc_register_service.html"
	loginIndex    = "https://m.myweimai.com/account/login.html?loginType=password"
	bookIndex     = "https://m.myweimai.com/wx/yb_book_new.html?from=doctor"
	dockerID      = "1483650534593388545" // 九价疫苗的
	// dockerID = "1367092118308483072" // 普通门诊的url

)

var (
	name         = "阮佩悦"
	myPhone      = "15267040566"
	seleniumPath = "/usr/local/bin/chromedriver"
	password     = "123456"
)

func findOrderIcon(wd selenium.WebDriver) (bool, error) {
	elements, err := wd.FindElements(selenium.ByCSSSelector, ".cir--base")
	if err != nil {
		return false, err
	} else if len(elements) > 0 {
		return true, err
	}
	return false, err
}

func findActiveOrderIcon(wd selenium.WebDriver) (bool, error) {
	elements, err := wd.FindElements(selenium.ByCSSSelector, ".cir--base.cir--active")
	if err != nil {
		return false, err
	} else if len(elements) > 0 {
		return true, err
	}
	return false, err
}

func captchaUnlock(wd selenium.WebDriver) (bool, error) {
	elem, err := wd.FindElement(selenium.ByCSSSelector, ".field-reqCode-unlock")
	if err != nil {
		return false, err
	} else if elem != nil {
		return true, err
	} else {
		return false, err
	}

}

func loginOk(wd selenium.WebDriver) (bool, error) {
	currentURL, err := wd.CurrentURL()
	if err != nil {
		return false, err
	} else if !strings.HasPrefix(currentURL, loginIndex) {
		return true, err
	} else {
		return false, err
	}

}

func pwdLogin(wd selenium.WebDriver) error {
	fmt.Println("页面需要登陆")
	phoneElem, err := wd.FindElement(selenium.ByCSSSelector, ".phone")

	if err != nil {
		return nil
	}
	phoneElem.Clear()
	err = phoneElem.SendKeys(myPhone)
	if err != nil {
		return err
	}

	pwdElem, err := wd.FindElement(selenium.ByCSSSelector, ".password")
	if err != nil {
		return err
	}
	if err := pwdElem.Clear(); err != nil {
		return err
	}

	err = pwdElem.SendKeys(password)
	if err != nil {
		return err
	}
	fmt.Println("数据密码登陆")
	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, ".login")
	if err != nil {
		return err
	}
	err = loginBtn.Click()
	if err != nil {
		return err
	}
	if err = wd.Wait(loginOk); err != nil {
		return err
	}

	return nil
}

func login(wd selenium.WebDriver) error {
	fmt.Println("页面需要登陆")
	phoneElem, err := wd.FindElement(selenium.ByCSSSelector, ".phone")

	if err != nil {
		return nil
	}
	phoneElem.Clear()
	err = phoneElem.SendKeys(myPhone)
	if err != nil {
		return err
	}
	codeElem, err := wd.FindElement(selenium.ByCSSSelector, ".code")
	if err != nil {
		return err
	}
	if err := codeElem.Clear(); err != nil {
		return err
	}

	captchaBtn, err := wd.FindElement(selenium.ByCSSSelector, ".field-reqCode")
	if err != nil {
		return err
	}
	// 等待可以发送激活码的时候再重新发
	if err := wd.Wait(captchaUnlock); err != nil {
		return err
	}

	captchaBtn.Click()
	fmt.Println("验证码已经发送")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("验证码已经发送到手机中, 并按enter登陆: ")
	captcha, _ := reader.ReadString('\n')
	captcha = strings.TrimSpace(captcha)
	fmt.Println(captcha)
	if err := codeElem.SendKeys(captcha); err != nil {
		return err
	}
	fmt.Println("填入验证码程序开始登陆")
	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, ".login")
	if err != nil {
		return err
	}
	err = loginBtn.Click()
	if err != nil {
		return err
	}

	return nil
}

func loadPatientOk(wd selenium.WebDriver) (bool, error) {
	elemt, err := wd.FindElement(selenium.ByCSSSelector, ".list-group.book-info")
	if err != nil {
		return false, err
	}
	elemts, err := elemt.FindElements(selenium.ByCSSSelector, ".list-item")

	if err != nil {
		return false, err
	}

	txt, err := elemts[0].Text()

	if err != nil {
		return false, err
	}

	if strings.Contains(txt, "请选择家庭成员") {
		return true, err
	}

	return true, err
}

func patientCardOk(wd selenium.WebDriver) (bool, error) {
	elemt, err := wd.FindElement(selenium.ByCSSSelector, ".list-group.book-info")
	if err != nil {
		return false, err
	}
	elemts, err := elemt.FindElements(selenium.ByCSSSelector, ".list-item")

	if err != nil {
		return false, err
	}

	elemt, err = elemts[1].FindElement(selenium.ByCSSSelector, ".right")
	if err != nil {
		return false, err
	}
	txt, err := elemt.Text()
	if err != nil {
		return false, err
	}
	fmt.Println(txt)
	if strings.HasPrefix(txt, "就诊卡") {
		return true, err
	}
	return false, err
}
func inCashURL(wd selenium.WebDriver) (bool, error) {
	url, err := wd.CurrentURL()
	if err != nil {
		return false, err
	}
	if strings.Contains(url, "cashier") {
		return true, nil
	}
	return false, nil
}

func payBtnOK(wd selenium.WebDriver) (bool, error) {
	_, err := wd.FindElement(selenium.ByClassName, "#addBtn2 > button")
	if err != nil {
		return false, err
	}
	return true, nil
}

func timeCardOk(wd selenium.WebDriver) (bool, error) {
	elemt, err := wd.FindElement(selenium.ByCSSSelector, ".list-group.book-time")
	if err != nil {
		return false, err
	}
	elemts, err := elemt.FindElements(selenium.ByCSSSelector, ".list-item")

	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}

	var elem selenium.WebElement
	for _, e := range elemts {
		txt, _ := e.Text()
		if strings.Contains(txt, "就诊号") {
			elem = e
		}
	}
	if elem != nil {
		btn, err := elem.FindElement(selenium.ByCSSSelector, ".right")
		if err != nil {
			return false, err
		}
		txt, err := btn.Text()
		if err != nil {
			return false, err
		}
		if txt == "请选择预约时段" {
			return false, nil
		}
	}
	return true, nil
}

func chooseTime(wd selenium.WebDriver) error {
	fmt.Println("开始就诊时间")
	if err := wd.WaitWithTimeout(timeCardOk, time.Millisecond*500); err != nil {
		wd.Refresh()
		chooseTime(wd)
	}
	return nil
}

func choosePaient(wd selenium.WebDriver) error {
	fmt.Println("开始选择病人")
	elemt, err := wd.FindElement(selenium.ByCSSSelector, ".list-group.book-info")
	if err != nil {
		return err
	}
	if err := wd.Wait(loadPatientOk); err != nil {
		return nil
	}

	elems, err := elemt.FindElements(selenium.ByCSSSelector, ".list-item")
	if err != nil {
		return err
	}

	// 0 就诊人， 1 就诊卡 ， 2 就诊费用，3 支付方式
	patientBtn, err := elems[0].FindElement(selenium.ByCSSSelector, ".right")
	if err != nil {
		return err
	}
	err = patientBtn.Click()
	if err != nil {
		return err
	}

	panelBody, err := wd.FindElement(selenium.ByClassName, "panel-body")
	if err != nil {
		return err
	}

	elemts, err := panelBody.FindElements(selenium.ByTagName, "li")
	if err != nil {
		return err
	}

	for _, e := range elemts {
		txt, _ := e.Text()
		txt = strings.TrimSpace(txt)
		if txt == name {
			fmt.Println("你选择的病人为: ", txt)
			return e.Click()
		}
	}
	fmt.Println("你选择的病人为当一个")
	return elemts[0].Click()
}

var exitCh = make(chan struct{})

func main() {

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

	go func() {
		select {
		case <-killSignal:
			fmt.Println("exit the code")
			os.Exit(0)
		}
	}()

	// opts := []selenium.ServiceOption{
	// 	selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	// 	selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	// 	selenium.Output(os.Stderr),            // Output debug information to STDERR.
	// }
	opts := []selenium.ServiceOption{}
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	imagCaps := map[string]interface{}{}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	defer wd.Close()

	// Navigate to the simple playground interface.
	encodedURL := url.PathEscape(fmt.Sprintf("%s?doctorId=%s", registerIndex, dockerID))
	loginURL := fmt.Sprintf("%s&redirect=%s", loginIndex, encodedURL)
	ch2 := make(chan int)

	go func(ch <-chan int) {
		<-ch

		for {
			if err := wd.WaitWithTimeout(findActiveOrderIcon, time.Second); err == nil {
				elements, err := wd.FindElements(selenium.ByCSSSelector, ".cir--base.cir--active")

				if err != nil {
					continue
				}
				if len(elements) == 0 {
					wd.Refresh()
					time.Sleep(time.Millisecond)
					continue
				}

				elem := elements[0]

				if err := elem.Click(); err != nil {
					continue
				}

				currentURL, err := wd.CurrentURL()

				if err != nil {
					continue
				}

				fmt.Println("当前网址为: ", currentURL)

				if currentURL == bookIndex {

					if err := chooseTime(wd); err != nil {
						fmt.Println("选择时间错误", err)
						continue
					}

					if err := choosePaient(wd); err != nil {
						continue
					}
				}
				if err = wd.Wait(patientCardOk); err != nil {
					continue
				}
				if err := chooseTime(wd); err != nil {
					continue
				}

				btn, err := wd.FindElement(selenium.ByCSSSelector, ".btn-book")
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = btn.Click()
				// 全部成功，break
				if err != nil {
					fmt.Println(err)
					continue
				}

				if err = wd.Wait(inCashURL); err != nil {
					continue
				}

				if err = wd.Wait(payBtnOK); err != nil {
					continue
				}

				payBtn, err := wd.FindElement(selenium.ByCSSSelector, "#addBtn2 > button")
				if err != nil {
					fmt.Println(err)
					continue
				}
				if txt, err := payBtn.Text(); err != nil {
					continue
				} else {
					fmt.Println(txt)
				}
				payBtn.Click()
				time.Sleep(timeDelta)
				fmt.Println("全部ok")
				break
			}
		}
	}(ch2)

	if err := wd.Get(loginURL); err != nil {
		panic(err)
	}

	currentURL, err := wd.CurrentURL()
	if err != nil {
		panic(err)
	}

	for {
		if !strings.HasPrefix(currentURL, loginIndex) {
			break
		}
		err = pwdLogin(wd)
		if err != nil {
			fmt.Printf("登陆失败，在%s后重试\r\n", timeDelta)
		}
		if newURL, err := wd.CurrentURL(); err != nil {
			fmt.Printf("登陆失败，在%s后重试\r\n", timeDelta)
		} else {
			currentURL = newURL
		}
		wd.Refresh()
	}

	// Waiting until login
	fmt.Println("没有找到order的按钮", time.Millisecond, "后重试")

	for {
		if err := wd.WaitWithTimeout(findActiveOrderIcon, time.Second); err != nil {
			fmt.Println("没有找到order的按钮", time.Millisecond, "后重试")
			wd.Refresh()
			time.Sleep(time.Millisecond)
		} else {
			break
		}
	}

	ch2 <- 1

	// 一把死循环各种调用各种跑

	// this channel forever only for block the main go routine stop
	<-exitCh
	fmt.Println("退出再见，拜拜")
}
