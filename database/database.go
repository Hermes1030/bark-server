package database

// Database defines all of the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	DeleteDeviceByKey(key string) error                     //Delete specified device
	Close() error                                           //Close the database
}

// AndroidDevice represents an Android device registration
type AndroidDevice struct {
    ID          string `json:"id"`
    DeviceToken string `json:"device_token"`
    DeviceName  string `json:"device_name"`
    Platform    string `json:"platform"` // "android"
    UserID      string `json:"user_id"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}

type DeviceManager interface {
    SaveDevice(device *AndroidDevice) error
    GetDevice(id string) (*AndroidDevice, error)
    DeleteDevice(id string) error
    // 其他方法...
}
