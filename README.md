
### JTRAN: JSON Transformation Framework

#### Overview

JTRAN is a powerful JSON transformation framework designed to help developers and organizations manipulate and transform JSON data according to predefined rules. This framework allows you to define transformation schemas or templates that dictate how data should be structured, processed, and presented. JTRAN is particularly useful for normalizing data from different sources, ensuring consistent output across various systems, such as APIs, front-end applications, or data migration tasks.

#### Key Features

1.  **Templating System**:
    
    -   **Dynamic Field References**: Use placeholders like  `${FIELD_NAME}`  within your JSON schema to dynamically reference fields from the input data.
    -   **Value Concatenation**: Combine multiple fields or values into a single output using templating, enabling complex transformations with minimal effort.
2.  **Keyword-Driven Transformation**:
    
    -   **Predefined Keywords**: JTRAN supports a set of built-in keywords that can manipulate data during transformation, such as  `uppercase`,  `capitalize`,  `replace`, and  `redact`.
    -   **Custom Keywords**: Extend the framework by registering your own custom keywords to introduce specialized transformation logic tailored to your specific use case.
3.  **Global and Scoped Transformations**:
    
    -   **Global Transformations**: Apply transformations that affect the entire JSON object.
    -   **Scoped Transformations**: Apply transformations that target specific fields or nested structures within the JSON object.
4.  **Array and Object Handling**:
    
    -   **Nested Structures**: JTRAN supports complex JSON structures, including arrays and nested objects, enabling deep transformations that can handle hierarchical data.
    -   **Array Operations**: Keywords like  `join`,  `foreach`, and others allow you to manipulate arrays, iterate over them, or apply transformations to each element.
5.  **Customizable and Extensible**:
    
    -   **Custom Keyword Handlers**: Register custom handlers to introduce new transformation logic or override existing ones, making the system highly adaptable to different business requirements.

#### Example Scenarios

1.  **Basic Transformation**: Suppose you have the following input JSON and want to transform it using JTRAN:
    
    **Input JSON**:
    
    json
    
    `{
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com"
    }` 
    
    **Schema**:
    
    json
    
    `{
        "full_name": "${first_name} ${last_name|uppercase}",
        "contact_email": "${email}"
    }` 
    
    **Transformed Output**:
    
    json
    
    `{
        "full_name": "John DOE",
        "contact_email": "john.doe@example.com"
    }` 
    
    **Explanation**:
    
    -   The  `full_name`  field is created by concatenating  `first_name`  and  `last_name`, with the  `last_name`  being transformed to uppercase.
    -   The  `contact_email`  field simply maps the  `email`  field from the input JSON.
2.  **Array Handling**: Let's say you have an array of user objects and want to join their full names into a single string.
    
    **Input JSON**:
    
    json
    
    `{
        "users": [
            {"first_name": "John", "last_name": "Doe"},
            {"first_name": "Jane", "last_name": "Smith"},
            {"first_name": "Alice", "last_name": "Johnson"}
        ]
    }` 
    
    **Schema**:
    
    json
    
    `{
        "all_users": "${users|foreach:first_name} ${users|foreach:last_name|join:, }"
    }` 
    
    **Transformed Output**:
    
    json
    
    `{
        "all_users": "John Doe, Jane Smith, Alice Johnson"
    }` 
    
    **Explanation**:
    
    -   The  `foreach`  keyword is used to iterate over the  `users`  array and extract the  `first_name`  and  `last_name`  of each user.
    -   The  `join`  keyword then concatenates these names into a single string, separated by commas.
3.  **Custom Keyword Example**: Imagine you need to mask sensitive data, such as credit card numbers, but still display the last four digits.
    
    **Input JSON**:
    
    json
    
    `{
        "credit_card": "1234-5678-9876-5432"
    }` 
    
    **Schema**:
    
    json
    
    `{
        "masked_credit_card": "${credit_card|redact(0,12)}${credit_card|substring(12,16)}"
    }` 
    
    **Transformed Output**:
    
    json
    
    `{
        "masked_credit_card": "****-****-****-5432"
    }` 
    
    **Explanation**:
    
    -   The  `redact`  keyword replaces the first 12 characters of the credit card number with asterisks.
    -   The  `substring`  keyword then extracts the last four digits of the credit card number to be appended to the masked result.

#### How It Works

1.  **Defining a Schema**: A schema is a JSON object that defines how the input data should be transformed. The schema may include simple field mappings, concatenation of fields, or more complex transformations using keywords.
    
    **Example Schema**:
    
    json
    
    `{
        "full_name": "${first_name} ${last_name|capitalize}",
        "email": "${email}"
    }` 
    
2.  **Applying Keywords**: Keywords are applied within the schema to transform the data. For example,  `${last_name|capitalize}`  capitalizes the  `last_name`  field, ensuring a consistent format.
    
3.  **Transforming Data**: The  `TransformData`  function takes the schema and the input JSON as arguments, processes the transformation rules, and returns the transformed JSON object.
    

#### Extending JTRAN

Developers can extend JTRAN by registering custom keywords using the  `RegisterKeyword`  function. This allows for the creation of specialized transformation logic that is not covered by the built-in keywords.

**Example: Registering a Custom Keyword**:

go

`RegisterKeyword("reverse", func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
    strVal, ok := value.(string)
    if !ok {
        return nil, fmt.Errorf("reverse keyword expects a string value")
    }
    
    runes := []rune(strVal)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }

    return string(runes), nil
})` 

In this example, the  `reverse`  keyword reverses the characters in a string. Once registered, it can be used within the schema like any other keyword.

#### Conclusion

JTRAN is a robust framework designed to simplify the process of transforming JSON data. Its combination of templating, keyword-driven transformations, and extensibility makes it a valuable tool for developers working with complex data structures in APIs, data migrations, or dynamic content generation. By defining transformation rules through schemas, JTRAN allows for consistent, maintainable, and adaptable data processing across various applications.