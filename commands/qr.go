package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/ozgur-yalcin/mfa/internal/database"
	"github.com/ozgur-yalcin/mfa/internal/initialize"
	"github.com/ozgur-yalcin/mfa/internal/models"
)

type qrCommand struct {
	r        *rootCommand
	fs       *flag.FlagSet
	commands []Commander
	name     string
	mode     string
	hash     string
	digits   int
	period   int64
	counter  int64
}

func newQrCommand() *qrCommand {
	qrCmd := &qrCommand{name: "qr"}
	return qrCmd
}

func (c *qrCommand) Name() string {
	return c.name
}

func (c *qrCommand) Commands() []Commander {
	return c.commands
}

func (c *qrCommand) Init(cd *Ancestor) error {
	c.fs = flag.NewFlagSet(c.name, flag.ExitOnError)
	c.fs.StringVar(&c.mode, "mode", "totp", "use time-variant TOTP mode or use event-based HOTP mode")
	c.fs.StringVar(&c.mode, "m", "totp", "use time-variant TOTP mode or use event-based HOTP mode (shorthand)")
	c.fs.StringVar(&c.hash, "hash", "SHA1", "A cryptographic hash method H")
	c.fs.StringVar(&c.hash, "H", "SHA1", "A cryptographic hash method H (shorthand)")
	c.fs.IntVar(&c.digits, "digits", 6, "A HOTP value digits d")
	c.fs.IntVar(&c.digits, "l", 6, "A HOTP value digits d (shorthand)")
	c.fs.Int64Var(&c.counter, "counter", 0, "used for HOTP, A counter C, which counts the number of iterations")
	c.fs.Int64Var(&c.counter, "c", 0, "used for HOTP, A counter C, which counts the number of iterations (shorthand)")
	c.fs.Int64Var(&c.period, "period", 30, "used for TOTP, an period (Tx) which will be used to calculate the value of the counter CT")
	c.fs.Int64Var(&c.period, "i", 30, "used for TOTP, an period (Tx) which will be used to calculate the value of the counter CT (shorthand)")
	return nil
}

func (c *qrCommand) Run(ctx context.Context, cd *Ancestor, args []string) error {
	initialize.Init()
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	imagePath := c.fs.Arg(0)
	otpauth, err := c.readQRCode(imagePath)
	if err != nil {
		return err
	}
	u, err := url.Parse(otpauth.String())
	if err != nil {
		return err
	}
	if u.Scheme != "otpauth" {
		return errors.New("invalid scheme")
	}
	account := &models.Account{}
	account.Mode = c.mode
	account.Hash = c.hash
	account.Digits = c.digits
	account.Period = c.period
	account.Counter = c.counter
	q := u.Query()
	if host := u.Hostname(); host != "" {
		account.Mode = host
	}
	if user := u.User.Username(); user != "" {
		account.User = user
	}
	if secret := q.Get("secret"); secret != "" {
		account.Secret = secret
	}
	if hash := q.Get("algorithm"); hash != "" {
		account.Hash = hash
	}
	if issuer := q.Get("issuer"); issuer != "" {
		account.Issuer = issuer
	}
	if digits := q.Get("digits"); digits != "" {
		fmt.Sscanf(digits, "%d", &account.Digits)
	}
	if period := q.Get("period"); period != "" && account.Mode == "totp" {
		fmt.Sscanf(period, "%d", &account.Period)
	}
	if counter := q.Get("counter"); counter != "" && account.Mode == "hotp" {
		fmt.Sscanf(counter, "%d", &account.Counter)
	}
	if err := c.saveAccount(account.Issuer, account.User, account.Secret, account.Mode, account.Hash, account.Digits, account.Period, account.Counter); err != nil {
		log.Fatal(err)
	}
	log.Println("account added successfully")
	return nil
}

func (c *qrCommand) readQRCode(imagePath string) (*gozxing.Result, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var img gozxing.LuminanceSource
	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".jpg", ".jpeg":
		img_obj, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		img = gozxing.NewLuminanceSourceFromImage(img_obj)
	case ".png":
		img_obj, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		img = gozxing.NewLuminanceSourceFromImage(img_obj)
	default:
		return nil, errors.New("unsupported image format")
	}
	bmp, err := gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(img))
	if err != nil {
		return nil, err
	}
	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *qrCommand) saveAccount(issuer string, user string, secret string, mode string, hash string, digits int, period int64, counter int64) error {
	db, err := database.LoadDatabase()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	accounts, err := db.ListAccounts(issuer, user)
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts) > 0 {
		log.Fatal("account already exists")
	} else if len(accounts) == 0 {
		account := &models.Account{
			Issuer:  issuer,
			User:    user,
			Secret:  secret,
			Mode:    mode,
			Hash:    hash,
			Digits:  digits,
			Period:  period,
			Counter: counter,
		}
		return db.CreateAccount(account)
	}
	return nil
}
