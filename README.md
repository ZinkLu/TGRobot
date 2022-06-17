# My TG Robot

Simplified Chinese document is [here](https://lh-love.top/posts/technical/%E8%87%AA%E7%94%A8tg%E6%9C%BA%E5%99%A8%E4%BA%BA/)

# 1. Introduction

a simple telegram robot that helps (especially with [vmshell server](https://vmshell.com/))
# 2. Build

go 1.18 is required

```
git clone git@github.com:ZinkLu/TGRobot.git && cd TGRobot && go build
```

# 3. Usage

## 3.1 Startup

to start the robot, a proper config file is required，check [config](###config) part to get more details.

command:

```
./TGRobot -c config.yaml
```
## 3.2 Config

a config file is a valid yaml format document or a json format document.

> yaml can have comments, which I suggest.

```yaml
apiToken: xxx
debug: false
handlers:
    message_handler:
        vmShell: xxx
        anotherMessageHandler: xxx
    inline_keyboard_handler:
        xxx: Xxx
    picture_handler:
        xxx: xxx
    ...
```

- `apiToken`: robot token which can be get from [@botfather](https://t.me/botfather)

- `debug`: verbose message will be logged if set to true.

- `handlers` section contains different handlers configs.

there are some handlers which can be used out of box:

## 3.3 VmShell Handler

vmshell handler helps you to get your vmshell server info or control your server for conveniently.

### 3.3.1 config

vmshell is a message handler yet a inline-keyboard handler，but since it process raw dialog messages, let's put it in `message_handler`: 

```yaml
handlers:
    message_handler:
        vmShell:
            username: vmshellAccount
            password: vmshellAccountPassword
            serverIds:
                - xxx
                - xxx
```

- `handers.message_handler.vmShell`:
    - `username`: vmshell account
    - `password`: vmshell password
    - `serverIds`: servers which you wants to control.

> warning!
> 
> a two step Authenticator should not be activated! 
> 
> until vmshell servers can be access through apiToken which is under developed according to their customer service.

#### 3.3.2 HOW TO GET serverId

1. open your services list

2. press `F12` to open develop console

3. remember to select `preserve log` and filter `Fetch/XHR`

4. select a server, just like the picture below:

    ![s1](docs/static/step1.jpg)

5. then your console should trace a XHR request which contains `serverId`:

    ![s2](docs/static/step2.jpg)


### 3.3.2 usage

currently, valid messages are:

- `流量`: get server bandwidth usage

- `服务器`: get server info

just send any message with theses keyword above to robot, it will retrieve the information for you.

### 3.3.3 TODO

- [x] make serverId to serverIds so we can control multiple servers.

- [ ] if serverIds is empty then robot can get all servers automatically for you to select.

## 3.4 Hitokoto / yiyan Handler

hitokoto handler doesn't need any configuration.

### 3.4.1 Usage

send `一句话` to robot to get your hitokoto.

thanks to [hitokoto.cn](https://hitokoto.cn/)

# 4. Add Custom Handler

## 4.1 Project Layout

since telegram have many message types, the source codes are structured to handle different types of message, we call the true handler `App Handler`.

```
├── handlers
│   ├── handlers.go
|   ├── register.go
│   └── message_handler
│       ├── message_handler.go
│       └── vmshell
│           ├── config.go
│           ├── server_info.go
│           ├── vmshell_client.go
│           ├── vmshell_client_test.go
│           └── vmshell_handler.go # this is a App Handler
|   ├── inline_keyboard_handler
|   ├── video_handler(not implement)
|   ├── command_handler(not implement)
```

## 4.2 Add App Handler

Assume you want to add a message handler to get the local weather.

First,  create a `weather` folder under `handlers/message_handler/`

## 4.3 Define Config

We want to specify a country in the config file so that we could get the city's weather.

let's add a mapping under `handler.message_handler` section in `config.yaml`

```yaml
apiToken: xxx
debug: false
handlers:
    message_handler:
        weather: # APP Handler Name
            city: Shanghai
```

Config under `message_handler` will be injected to App Handler automatically.

## 4.4 Define Handler

An App Handler must implements `common.AppHandlerInterface`.

```golang
type AppHandlerInterface interface {
	Handle(*tgbotapi.Update, *tgbotapi.BotAPI) // Handler function
	When(*tgbotapi.Update) bool // true means the handler can handler current message, or fallback next handler
	Init(*config.ConfigUnmarshaler) // Init function can be called automatically with config as it's parameters
	Order() int // less is earlier 
	Help() string // Help String, If all App can't handler current message, a combination of help messages is sent by bot 
	Name() string // Name of the App Handler , can't duplicate.
}
```

In order to use the configuration file in the yaml, we define a struct that corresponds to it.

```golang
package weather

type Config struct {
	City string `configKey:"city"`
}
```

Since we can take json and yaml as our config file, a new struct tag named `configKey` is used to unmarshal a config object.

```golang
func (w *WeatherHandler) Init(conf *config.ConfigUnmarshaler) {
	wConf := &Config{}
	conf.UnmarshalConfig(wConf, w.Name())
	w.City = wConf.City
}

func (w *WeatherHandler) Name() string {
	return "weather"
}

func (w *WeatherHandler) Order() int {
	return 999
}
```

`*config.ConfigUnmarshaler`'s `UnmarshalConfig` method will pass the config under `handler.message_handler` section for you.

By the way, we set `Order()` function, This method also affects the order of help messages.

Let's set `When()` can this App Handler handler messages.

Let's say, the handler can handle messages which contains word "weather".

```golang
func (w *WeatherHandler) When(u *tgbotapi.Update) bool {
	return strings.Contains(u.Message.Text, "weather")
}

func (w *WeatherHandler) Help() string {
    return "ask me the 'weather'"
}
```

> `u.Message` is a pointer, but you should not worry it will be a nil, cause message handler only can handle a non-nil `u.Message`

Now that the `Handle()` method can be implements , let's write a pseudo-code

```golang
func (w *WeatherHandler) Handle(u *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	URL, _ := url.Parse(WEATHER_PROVIDER_URL) 
	URL.Query().Set("City", w.City)
	resp, _ := http.Get(URL.String())
	content, _ := ioutil.ReadAll(resp.Body)
	m := tgbotapi.NewMessage(u.Message.Chat.ID, string(content))
	bot.Send(m)
}
```

## 4.5 Register Handler

In order to enable the Handler, you need to register the Handler to its parent Handler, our weather handler's parent handler is `message_handler`, call `message_handler.Register` directly in `init`.

```golang
func init() {
	message_handler.Register(&WeatherHandler{})
}
```

So that we can import the package in `handlers.go` to enable the handler.

```golang
package handlers

import (
    _ "github.com/ZinkLu/TGRobot/handlers/message_handler/weather"
)
```

## 4.6 Get other handler

Sometimes an App Handler can depend on other handlers.

for example an `inline_keyboard_message` may be triggered by a message handler.

So a `inline_keyboard_message handler` may need the trigger's config to process further.

with `GetAppHandlerByName()` function, you can get a handler that has been registered.

if the handler is not registered, the program may panic.

```golang
message_handler := pool.GetAppHandlerByName[*vm_message.VmShellHandler]("vmShell")
message_handler.Config.serverIds // get info
```
