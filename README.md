# Apspace-calendar
Automatically sync your APSpace timetable with most calendar apps using an iCalendar URL.

An instance is hosted at [apspace-calendar.netlify.app](https://apspace-calendar.netlify.app/)

## Usage
### 1. Getting the url
Example: `https://apspace-calendar.netlify.app/?intake=AFCF2411ICT&group=G1&title=module_name`

#### Parameters
- `intake`
- `group` (optional) - your group, e.g. "G1". if you don't specify anything, classes from all groups will be given.
    - Do tell me if there are any issues with this feature in particular
- `title` (optional) - what to use for the calendar event title
    - `module_name` (default) - e.g., **Basic Marketing Skills**.
    - `module_name_class` - e.g., **Basic Marketing Skills T-1**.
    - `module_code` - e.g., **BMS**.
    - `module_code_class` - e.g., **BMS T-1**.
    - `module_id` - the full module ID, which is what APSpace's export feature uses, e.g., **ABUS012-4-C-BMS-T-1**.

### 2. Adding it to your calendar app
For example: Google Calendar

1. Go to calendar.google.com
2. Go to settings
3. Go to Add calendar > From URL
4. Enter your URL and click Add calendar

## Running it for yourself
### `net/http` version
```sh
go get ./...
go run src/main.go --serve
```

### Netlify function
Just deploy to netlify. You can also run it locally with `netlify dev`.

## Limitations/TODO
- [x] Filter by grouping
- [x] Tests
- [ ] Ability to filter classes
- [ ] urm.. Cache the data for the net/http version?

## Contributing
Feel free to open an issue for any bugs. Or better yet, fix it for me and open a pull request 🥺.
