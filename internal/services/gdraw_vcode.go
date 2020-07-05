package services

import (
	"github.com/mojocn/base64Captcha"
)

var GdrawVcodeer gdrawVcodeer

func init() {
	GdrawVcodeer = &gdrawVcode{
		Store: base64Captcha.DefaultMemStore,
		Mtp:   base64Captcha.DriverMath{Width: 200, Height: 100},
	}
}

type gdrawVcode struct {
	Store base64Captcha.Store
	Mtp   base64Captcha.DriverMath
}

type gdrawVcodeer interface {
	Get() (id, b64s string, err error)
	Verify(id string, answer string) (match bool)
}

func (this *gdrawVcode) Get() (id, b64s string, err error) {
	driver := this.Mtp.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, this.Store)
	id, b64s, err = c.Generate()
	return
}

func (this *gdrawVcode) Verify(id string, answer string) (match bool) {
	if id == "" || answer == "" {
		return
	}
	driver := this.Mtp.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, this.Store)
	return c.Verify(id, answer, true)
}
