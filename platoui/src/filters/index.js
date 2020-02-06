import Vue from 'vue'

Vue.filter('secondsAsDuration', function(seconds) {
  let hoursOnly = Math.floor(seconds / 3600)
  let minutesOnly = Math.floor((seconds % 3600) / 60)
  let secondsOnly = Math.floor(seconds % 60)
  let durationString = secondsOnly + "s"
  if (hoursOnly > 0 || minutesOnly > 0) {
    durationString = minutesOnly + "m " + durationString
  }
  if (hoursOnly > 0) {
    durationString = hoursOnly + "h " + durationString
  }
  return durationString
})
