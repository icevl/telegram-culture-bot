## Overview

Telegram bot with channels and prompt config.

```
[
 {
  "title": "Flow example #1",
  "chat_id": -1002062592133,
  "prompt": "Generate an interesting fact about any country and its culture, and not just about Japan. At the end of the text, write the name of the country after a '|' separator, followed by another '|' separator and an emoji for the country",
  "min_mins": 120, // Random interval from
  "max_mins": 180, // Random interval to
  "image": "from_prompt_result",
  "next_time": 1703096445 // Next time for bot sending
 },
 {
  "title": "Flow example #2",
  "chat_id": -1002062592133,
  "prompt": "Generate an interesting piece of wisdom from any country. At the end of the text, write 'As they say in X' after a '|' separator, where X is the country that this wisdom came from, followed by another '|' separator and an emoji for the country",
  "min_mins": 60,
  "max_mins": 120,
  "image": "from_prompt_result",
  "next_time": 1703093313
 }
]
```

image - DALLE-3 image generation. Can be:

- blank - without image
- "from_prompt_result" - result from GPT prompt goes to DALLE
- static text  

Example: [https://t.me/cultureandcountries](https://t.me/cultureandcountries)https://t.me/cultureandcountries
