import axios from 'axios'

axios.defaults.baseURL = '/api'

const getTracks = function() {
  return axios.get("/tracks")
}

export default {
  getTracks
}
