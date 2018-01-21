/* global localStorage */

import axios from 'axios'

const API_URL = process.env.API_URL || 'https://tips.mycoralhealth.com/'

export default axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})
