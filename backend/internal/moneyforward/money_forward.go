package moneyforward

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	loginURL    = "https://moneyforward.com/users/sign_in"
	downloadURL = "https://moneyforward.com/cf/csv"
)

// MoneyForward scpaes from Money Forward ME.
type MoneyForward struct {
	email    string
	password string
	cookies  []*http.Cookie
}

// New returns a new MoneyForward instance.
func New(email, password string) *MoneyForward {
	return &MoneyForward{
		email:    email,
		password: password,
	}
}

// Login does actions for login and get cookies.
//
// - GET /users/sign_in
// 	- redirect twice
// - GET メールアドレスでログイン
// - POST /sign_in/email
// 	- redirect to password
// - GET /password
// - POST /sign_in
// 	- redirect three times
// 	- finally redirect to https://moneyforward.com/
func (m *MoneyForward) Login(ctx context.Context) error {
	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate(loginURL),
	); err != nil {
		return fmt.Errorf("failed to get login url: %w", err)
	}

	// メールアドレスでログインをクリック
	loginWithEmailSelector := `a.ssoLink`
	if _, err := chromedp.RunResponse(chromeCtx,
		chromedp.Click(loginWithEmailSelector),
	); err != nil {
		return fmt.Errorf("failed to click login with email button: %w", err)
	}

	// メールアドレスを入力して送信
	emailSelector := `//input[@name="mfid_user[email]"]`
	if err := chromedp.Run(chromeCtx,
		chromedp.SendKeys(emailSelector, m.email),
		chromedp.Submit(emailSelector),
	); err != nil {
		return fmt.Errorf("failed to submit email: %w", err)
	}

	// パスワードを入力して送信
	passwordSelector := `//input[@name="mfid_user[password]"]`
	if _, err := chromedp.RunResponse(chromeCtx,
		chromedp.SendKeys(passwordSelector, m.password),
		chromedp.Submit(passwordSelector),
		chromedp.WaitNotPresent(passwordSelector),
	); err != nil {
		return fmt.Errorf("failed to submit password: %w", err)
	}

	// クッキーを取得
	var cookies []*network.Cookie
	if err := chromedp.Run(chromeCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			c, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			cookies = c
			return nil
		}),
	); err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	for _, c := range cookies {
		m.cookies = append(m.cookies, &http.Cookie{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	return nil
}

// DownloadCSV gets CSV given year and month.
//
// DownloadCSV returns CSV body as io.ReadCloser, so the caller must close it.
func (m *MoneyForward) DownloadCSV(ctx context.Context, year int, month time.Month) (io.ReadCloser, error) {
	v := url.Values{}
	v.Set("from", fmt.Sprintf("%04d/%02d/01", year, month))
	v.Set("month", fmt.Sprintf("%d", month))
	v.Set("year", fmt.Sprintf("%d", year))

	url := fmt.Sprintf("%s?%s", downloadURL, v.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new request: %w", err)
	}

	for _, c := range m.cookies {
		req.AddCookie(c)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return resp.Body, nil
}
