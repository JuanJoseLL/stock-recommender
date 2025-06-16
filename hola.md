To begin, you will need to connect to the API that provides the stock information. This will involve making HTTP requests to the API endpoints to retrieve the necessary data. You should handle errors appropriately and ensure that the data is properly formatted for use in the UI. After retrieving the data store it in CockroachDB. Make sure you download all the data from the API.
Your task is to design and implement a system that retrieves stock information from a given API. The system should be able to handle different types of stock and display the relevant information to the user in a user-friendly interface.
Endpoint

GET https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list
Query Parameters

Available parameters:

next_page: The key to start the next page
Authentication

Include API key in Authorization header:

Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MSwiZW1haWwiOiJqdWFuam9sbzEyQGhvdG1haWwuY29tIiwiZXhwIjoxNzUwMTg2MTk2LCJpZCI6IjAiLCJwYXNzd29yZCI6Iicgb3IgMT0xIG9yICcnPScifQ.i3iwFvk5BA5l_AonHDYb4RhEJSZfha3nVsb4ZFvQ7u0
{
    "items": [
        {
            "ticker": "BSBR",
            "target_from": "$4.20",
            "target_to": "$4.70",
            "company": "Banco Santander (Brasil)",
            "action": "upgraded by",
            "brokerage": "The Goldman Sachs Group",
            "rating_from": "Sell",
            "rating_to": "Neutral",
            "time": "2025-01-13T00:30:05.813548892Z"
        },
        {
            "ticker": "VYGR",
            "target_from": "$11.00",
            "target_to": "$9.00",
            "company": "Voyager Therapeutics",
            "action": "reiterated by",
            "brokerage": "Wedbush",
            "rating_from": "Outperform",
            "rating_to": "Outperform",
            "time": "2025-01-14T00:30:05.813548892Z"
        }
    ],
    "next_page": "AZEK"
}