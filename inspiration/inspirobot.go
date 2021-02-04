package inspiration

import (
	"io/ioutil"
	"net/http"
)

const (
	generate_url = "https://inspirobot.me/api?generate=true"
)

func NewInspirationalMessage() (string, error) {
	resp, err := http.Get(generate_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// class InspirobotClient:
//     """Client for interacting with inspirobot"""
//     def __init__(self):
//         self.generate_url = INSPIROBOT_GENERATE_URL

//     @staticmethod
//     async def fetch(session, url):
//         """Performs a GET on a URL"""
//         async with session.get(url) as response:
//             return await response.text()

//     async def generate_inspirational_message(self):
//         """Gets an inspirational message from inspirobot"""
//         async with aiohttp.ClientSession() as session:
//             html = await self.fetch(session, self.generate_url)
//             return html
