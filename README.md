# Apspace-calendar
Automatically sync your APSpace timetable with most calendar apps using an iCalendar URL.

An instance is hosted at [apspace-calendar.netlify.app](https://apspace-calendar.netlify.app/)

## Usage
### 1. Getting the url
Example: `https://apspace-calendar.netlify.app/?intake=AFCF2411ICT&title=module_name`

#### 2. Parameters
- `intake`
- `title` (optional) - what to use for the calendar event title
    - `module_name` (default): Uses the full module name, e.g., **Basic Marketing Skills**.
    - `module_code`: Extracts the module code from the module ID, e.g., **BMS**.
    - `module_id`: Uses the full module ID, which is what APSpace's export feature uses, e.g., **ABUS012-4-C-BMS-T-1**.

### Adding it to your calendar app
For example: Google Calendar

1. Go to calendar.google.com
2. Go to settings
3. Go to Add calendar > From URL
4. Enter your URL and click Add calendar

## Limitations/TODO
- No support for groupings

## Running it for yourself
### net/http version
```sh
go get ./...
go run src/main.go --serve
```

### netlify function
Just deploy to netlify. You can also run it locally with `netlify dev`.
