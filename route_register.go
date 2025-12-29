package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mritd/logger"
	"github.com/finb/bark-server/v2/database"
)

type DeviceInfo struct {
	DeviceKey   string `form:"device_key,omitempty" json:"device_key,omitempty" xml:"device_key,omitempty" query:"device_key,omitempty"`
	DeviceToken string `form:"device_token,omitempty" json:"device_token,omitempty" xml:"device_token,omitempty" query:"device_token,omitempty"`

	// compatible with old req
	OldDeviceKey   string `form:"key,omitempty" json:"key,omitempty" xml:"key,omitempty" query:"key,omitempty"`
	OldDeviceToken string `form:"devicetoken,omitempty" json:"devicetoken,omitempty" xml:"devicetoken,omitempty" query:"devicetoken,omitempty"`
}


func init() {
	registerRoute("register", func(router fiber.Router) {
		router.Post("/register", func(c *fiber.Ctx) error { return doRegister(c, false) })
		router.Get("/register/:device_key", doRegisterCheck)
	})

	// compatible with old requests
	registerRouteWithWeight("register_compat", 100, func(router fiber.Router) {
		router.Get("/register", func(c *fiber.Ctx) error { return doRegister(c, true) })
	})
	
	// Android device registration
	registerRouteWithWeight("android", 60, func(router fiber.Router) {
		router.Post("/register/android", registerAndroidDevice)
	})
}

func doRegister(c *fiber.Ctx, compat bool) error {
	var deviceInfo DeviceInfo
	if compat {
		if err := c.QueryParser(&deviceInfo); err != nil {
			return c.Status(400).JSON(failed(400, "request bind failed1: %v", err))
		}
	} else {
		if err := c.BodyParser(&deviceInfo); err != nil {
			return c.Status(400).JSON(failed(400, "request bind failed2: %v", err))
		}
	}

	if deviceInfo.DeviceKey == "" && deviceInfo.OldDeviceKey != "" {
		deviceInfo.DeviceKey = deviceInfo.OldDeviceKey
	}

	if deviceInfo.DeviceToken == "" {
		if deviceInfo.OldDeviceToken != "" {
			deviceInfo.DeviceToken = deviceInfo.OldDeviceToken
		} else {
			return c.Status(400).JSON(failed(400, "device token is empty"))
		}
	}

	// DeviceToken length is variable, but should not be too long.
	if len(deviceInfo.DeviceToken) > 128 {
		return c.Status(400).JSON(failed(400, "device token is invalid"))
	}

	// if deviceInfo.DeviceKey=="", newKey will be filled with a new uuid
	// otherwise it equal to deviceInfo.DeviceKey
	newKey, err := db.SaveDeviceTokenByKey(deviceInfo.DeviceKey, deviceInfo.DeviceToken)
	if err != nil {
		logger.Errorf("device registration failed: %v", err)
		return c.Status(500).JSON(failed(500, "device registration failed: %v", err))
	}
	deviceInfo.DeviceKey = newKey

	return c.Status(200).JSON(data(map[string]string{
		// compatible with old resp
		"key":          deviceInfo.DeviceKey,
		"device_key":   deviceInfo.DeviceKey,
		"device_token": deviceInfo.DeviceToken,
	}))
}

func doRegisterCheck(c *fiber.Ctx) error {
	deviceKey := c.Params("device_key")

	if deviceKey == "" {
		return c.Status(400).JSON(failed(400, "device key is empty"))
	}

	_, err := db.DeviceTokenByKey(deviceKey)
	if err != nil {
		return c.Status(400).JSON(failed(400, "%s", err.Error()))
	}
	return c.Status(200).JSON(success())
}

func registerAndroidDevice(c *fiber.Ctx) error {
    // 解析请求体
    var device database.AndroidDevice
    if err := c.BodyParser(&device); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // 验证必要字段
    if device.DeviceToken == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Device token is required"})
    }

    // 设置默认值
    if device.ID == "" {
        device.ID = generateDeviceID() // 使用现有工具函数
    }
    device.Platform = "android"
    device.CreatedAt = time.Now().Unix()
    device.UpdatedAt = time.Now().Unix()

    // 保存到数据库
    // TODO: 实现数据库层的SaveDevice方法
    // err := db.SaveDevice(&device) // 使用数据库层接口
    // if err != nil {
    //     return c.Status(500).JSON(fiber.Map{"error": "Failed to register device"})
    // }

    return c.JSON(fiber.Map{
        "success": true,
        "device_id": device.ID,
    })
}


