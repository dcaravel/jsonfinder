# JSON Finder

Utility for searching JSON documents.

Will output the matching values in the original JSON structure (with unmatched items removed). 

Optionally - unmatched fields (a.k.a. context) can be included in the output as well.

## Usage

```
jsonfinder -f <filename> "<search term>"
```

## Example

```json
{
    "A": {
        "label": "Label for A",
        "things": [
            "one two three",
            "two three four",
            "three four five"
        ]
    },
    "B": {
        "label": "Label for D",
        "things": [
            "hello",
            "world"
        ]
    },
    "C": {
        "label": "Label for C",
        "things": [
            "four",
            "five",
            "six"
        ],
        "bonus": [
            {
                "name": "dave",
                "favorite": "seven"
            },
            {
                "name": "fred",
                "favorite": "four"
            }
        ]
    }
}
```

```
$ jsonfinder -f example.json four
{
    "A": {
        "things": [
            "two three four",
            "three four five"
        ]
    }
    "C": {
        "things": [
            "four"
        ],
        "bonus": [
            {
                "favorite": "four"
            },
        ]
    }
}
```
To add context:
```
$ jsonfinder -f example.json --indexes --context "A label, C label, C bonus name" four
{
    "A": {
        "label":"Label for A",
        "things": [
            "two three four",
            "three four five"
        ]
    }
    "C": {
        "label":"Label for C",
        "things": [
            "four"
        ],
        "bonus": [
            {
                "_oindex": 1,
                "name": "fred",
                "favorite": "four"
            },
        ]
    }
}
```

## TODO

- Change 'context' to use jsonpath instead
- Allow regex in search term
- Change how breadcrumbs are tracked so that the output logic is as messy trying to interpret the before and after elements.
- Allow the field name for 'addindexes" to be customized
- Add support for searching bool and float fields, also ensure that these same field types can be included as context