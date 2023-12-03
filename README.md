# go-calendar

Generate a calendar from a JSON file.

## Options with args

| Option  | Description             | Default   |
| ------- | ----------------------- | --------- |
| -k KEY  | Key of the date         | date      |
| -c KEY  | Key of the counter      | num       |
| -i FILE | Path to the JSON file   | data.json |
| -o FILE | Path to the output file | out.svg   |
| -q      | Quiet mode              | false     |

## Usage

There are differents possible usages depending on your JSON file.

<details>

<summary>JSON file with dates </summary>

```json
[
  {
    "date": "2022-11-06",
  },
  {
    "date": "2022-11-07",
  },
]
```

</details>

You will need to specify the key of the date and the key of the counter

```sh
go run main.go -k date
```

<details>

<summary>JSON file with a counter</summary>

```json
[
  {
    "date": "2022-11-06",
    "num": 9
  },
  {
    "date": "2022-11-06",
    "num": 1
  },
]
```

</details>

You will need to specify the key of the date and the key of the counter

```sh
go run main.go -k date -c num
```

## Example

<details>

<summary>Generate calendar from Github contributions</summary>

Using gh cli and jq, you can get your contributions from Github with

```sh
gh api graphql -F owner='Its-Just-Nans' -f query='
    query( $owner: String!) {
      user(login: $owner) {
    contributionsCollection {
      contributionCalendar {
        totalContributions
        weeks {
          contributionDays {
            contributionCount
            weekday
            date
          }
        }
      }
    }
  }}' | jq '[.data.user.contributionsCollection.contributionCalendar.weeks | .[].contributionDays |.[] | {date: (.date), num:(.contributionCount)}]' > out.json
```

Then you can generate the calendar with

```sh
go run main.go -k date -c num -i out.json -o contributions.svg
```

</details


