curl -s --user 'api:key-d109e953a537e26356d3bf45a31ae065' \
    https://api.mailgun.net/v3/mail.mycoralhealth.com/messages \
    -F from='Health Tips <mailgun@mail.mycoralhealth.com>' \
    -F to=andy@mycoralhealth.com \
    -F subject='Hello here is a health tip' \
    -F text='Testing some Mailgun awesomeness!'
