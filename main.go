package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
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
	loginIndex    = "https://m.myweimai.com/account/login.html"
	bookIndex     = "https://m.myweimai.com/wx/yb_book_new.html?from=doctor"
	dockerID      = "1483650534593388545" // 九价疫苗的
	// dockerID = "1367092118308483072" // 普通门诊的url

)

var (
	myPhone, name, seleniumPath string
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
		return false, err
	}

	return true, err
}
func chooicePatient(wd selenium.WebDriver) error {

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

func checkName() {
	if name == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入姓名, 并按enter: ")
		yourName, _ := reader.ReadString('\n')
		name = strings.TrimSpace(yourName)
	}
}

func checkPhone() {
	if myPhone == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入手机号, 并按enter: ")
		phone, _ := reader.ReadString('\n')
		myPhone = strings.TrimSpace(phone)
	}
}

func checkWebDriver() {
	if seleniumPath == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入chromedriver路径, 并按enter: ")
		myPath, _ := reader.ReadString('\n')
		seleniumPath = strings.TrimSpace(myPath)
	}
	if _, err := os.Stat(seleniumPath); os.IsNotExist(err) {
		panic(err)
	}
}

func main() {
	checkName()
	checkPhone()
	checkWebDriver()
	// "/usr/local/bin/chromedriver"

	// opts := []selenium.ServiceOption{
	// 	selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	// 	selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	// 	selenium.Output(os.Stderr),            // Output debug information to STDERR.
	// }
	opts := []selenium.ServiceOption{}
	selenium.SetDebug(false)
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

	// Navigate to the simple playground interface.
	encodedURL := url.PathEscape(fmt.Sprintf("%s?doctorId=%s", registerIndex, dockerID))
	loginURL := fmt.Sprintf("%s?redirect=%s", loginIndex, encodedURL)

	if err := wd.Get(loginURL); err != nil {
		panic(err)
	}

	currentURL, err := wd.CurrentURL()
	if err != nil {
		panic(err)
	}

	for strings.HasPrefix(currentURL, loginIndex) {
		err = login(wd)
		if err != nil {
			fmt.Printf("登陆失败，在%s后重试\r\n", timeDelta)
			time.Sleep(timeDelta)
		}
		if newURL, err := wd.CurrentURL(); err != nil {
			fmt.Printf("登陆失败，在%s后重试\r\n", timeDelta)
			time.Sleep(timeDelta)
		} else {
			currentURL = newURL
		}
		wd.Refresh()
	}

	// Waiting until login
	if err := wd.Wait(findOrderIcon); err != nil {
		panic(err)
	}
	// 一把死循环各种调用各种跑
	for {
		elements, _ := wd.FindElements(selenium.ByCSSSelector, ".cir--base.cir--active")
		if len(elements) == 0 {
			fmt.Printf("没有可以预约的窗口，在%s后重试\r\n", timeDelta)
			time.Sleep(timeDelta)
			wd.Refresh()
		} else {
			elem := elements[0]
			txt, err := elem.Text()

			if err != nil {
				continue
			}

			fmt.Println("按钮的文本为: ", txt)
			// click the first book order

			if err := elem.Click(); err != nil {
				continue
			}

			currentURL, err = wd.CurrentURL()

			if err != nil {
				continue
			}

			fmt.Println("当前网址为: ", currentURL)

			if currentURL == bookIndex {
				if err := chooicePatient(wd); err != nil {
					continue
				}
			}
			// 全部成功，break
			break
		}
	}
	// this channel forever only for block the main go routine stop
}
