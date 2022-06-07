package constant

const (
    ProcessingState int = 1
    SuccessState    int = 2
    FailedState     int = 3
    AllState        int = 0
)

const (
    ModeTrimLeadingSpace bool   = true  // Ignore whitespace
    ModeComma            string = ":"   // Separator character
    ModeLazyQuotes       bool   = true // If LazyQuotes is true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field.
    ModeFieldsPerRecord  int    = 0     // Read sets it to the number of fields in the first record
)
