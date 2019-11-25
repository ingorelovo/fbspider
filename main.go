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

// Start a Selenium WebDriver server instance (if one is not already
// running).
// Start a Selenium WebDriver server instance (if one is not already
// running).
const (
	registerIndex = "https://m.myweimai.com/wx/dc_register_service.html"
	loginIndex    = "https://m.myweimai.com/account/login.html"
	bookIndex     = "https://m.myweimai.com/wx/yb_book_new.html?from=doctor"
	// dockerId = "1483650534593388545" // 九价疫苗的
	dockerID = "1367092118308483072" // 普通门诊的url
	myPhone  = "xxxxxxxx"            // 普通门诊的url
	name     = "陈超"
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
func chooicePatient(wd selenium.WebDriver) error {
	fmt.Println("开始选择病人")
	elemt, err := wd.FindElement(selenium.ByCSSSelector, ".list-group.book-info")
	if err != nil {
		return err
	}
	elems, err := elemt.FindElements(selenium.ByCSSSelector, ".list-item")
	if err != nil {
		return err
	}
	// 0 就诊人， 1 就诊卡 ， 2 就诊费用，3 支付方式
	if err := elems[0].Click(); err != nil {

	}
	return nil
}

func main() {

	const (
		// These paths will be different on your system.
		seleniumPath = "/usr/local/bin/chromedriver"
		port         = 8080
	)
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
			fmt.Println("retry after 0.1s")
			time.Sleep(time.Second * 5)
		}
		if newURL, err := wd.CurrentURL(); err != nil {
			fmt.Println("retry after 0.1s")
			time.Sleep(time.Second * 5)
		} else {
			currentURL = newURL
		}
		wd.Refresh()
	}

	// Waiting until login
	if err := wd.Wait(findOrderIcon); err != nil {
		panic(err)
	}
	// Get a order list
	elements, err := wd.FindElements(selenium.ByCSSSelector, ".cir--base.cir--active")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(elements))

	elem := elements[0]
	txt, err := elem.Text()
	if err != nil {
		panic(err)
	}
	fmt.Println("按钮的文本为: ", txt)
	// click the first book order
	if err := elem.Click(); err != nil {
		panic(err)
	}
	currentURL, err = wd.CurrentURL()
	if err != nil {
		panic(err)
	}
	fmt.Println("current url: ", currentURL)
	if currentURL == bookIndex {
		if err := chooicePatient(wd); err != nil {
			panic(err)
		}
	}
	// this channel forever only for block the main go routine stop
	forever := make(chan bool)

	<-forever
	// // Get a reference to the text box containing code.
	// elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
	// if err != nil {
	// 	panic(err)
	// }
	// // Remove the boilerplate code already in the text box.
	// if err := elem.Clear(); err != nil {
	// 	panic(err)
	// }

	// // Enter some new code in text box.
	// err = elem.SendKeys(`
	// 	package main
	// 	import "fmt"

	// 	func main() {
	// 		fmt.Println("Hello WebDriver!\n")
	// 	}
	// `)
	// if err != nil {
	// 	panic(err)
	// }

	// // Click the run button.
	// btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	// if err != nil {
	// 	panic(err)
	// }
	// if err := btn.Click(); err != nil {
	// 	panic(err)
	// }

	// // Wait for the program to finish running and get the output.
	// outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "#output")
	// if err != nil {
	// 	panic(err)
	// }

	// var output string
	// for {
	// 	output, err = outputDiv.Text()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if output != "Waiting for remote server..." {
	// 		break
	// 	}
	// 	time.Sleep(time.Millisecond * 100)
	// }

	// fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	// // Example Output:
	// // Hello WebDriver!
	// //
	// // Program exited.
}
