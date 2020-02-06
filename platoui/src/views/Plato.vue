<template>
  <div class="plato">
    <div>
      <p>Plato</p>
        <div>
          <button @click="backTrack">BACK</button>
          <button @click="playPause">PLAY/PAUSE</button>
          <button @click="nextTrack">NEXT</button>
        </div>
        <div>
          VOL <input v-model="volume" type="range" min="0" max="1.0" step="0.05">
        </div>
        <div>
          {{ seek | secondsAsDuration }} <input v-model="seek" type="range" min="0" :max="currentHowlLength" step="1"> {{ currentHowlLength | secondsAsDuration }}
        </div>
        <ul>
          <li v-for="item in playlist" :key="item.path" :class="{active: currentItem == item}">
            <template v-if="item.artist">{{ item.artist }} - </template>
            {{ item.title }}
          </li>
        </ul>
    </div>
    <pre>{{ tracks }}</pre>
  </div>
</template>

<script>
import { Howler, Howl } from 'howler'
import api from '@/tools/api'
import '@/filters'

export default {
  name: 'plato',
  data: function () {
    return {
      playing: false,
      playlist: [
        {title: "Song Title", path: "musics/testencoding/Beck - Loser.mp3"},
      ],
      currentIndex: 0,
      currentHowl: null,
      manuallyChangingSeek: false,
      seek: 0,
      tracks: [],
      volume: Howler.volume(),
    }
  },
  computed: {
    currentItem: function() {
      return this.playlist[this.currentIndex]
    },
    currentHowlLength: function() {
      if (this.currentHowl) {
        return this.currentHowl.duration()
      }
      return null
    },
  },
  watch: {
    playing: function(playing) {
      let vm = this
      vm.seek = vm.currentHowl.seek()
      let updateSeek
      if (playing) {
        updateSeek = setInterval(function() {
          vm.seek = vm.currentHowl.seek()
        }, 200)
      } else {
        clearInterval(updateSeek)
      }
    },
    tracks: function() {
      this.makePlaylistFromTracks()
    },
    volume: function(value) {
      Howler.volume(value)
    },
  },
  created: function() {
    this.fetchTracks()
    this.currentHowl = this.makeNewHowl(this.currentItem.path)
  },
  methods: {
    backTrack: function() {
      if (this.currentHowl) {
        this.currentHowl.stop()
      }
      this.currentIndex--
      if (this.currentIndex < 0) {
        this.currentIndex = this.playlist.length - 1
      }
      this.currentHowl = this.makeNewHowl(this.currentItem.path)
      if (this.playing) {
        this.currentHowl.play()
      }
    },
    fetchTracks: function() {
      var vm = this
      api.getTracks()
      .then(function(response) {
        vm.tracks = response.data
      })
      .catch(function(error) {
        console.error("Cannot retrieve tracks: ", error)
      })
    },
    makeNewHowl: function(filePath) {
        var sound = new Howl({
            src: filePath,
        });
        return sound
    },
    makePlaylistFromTracks: function() {
      this.playlist = []
      var vm = this
      for (const track of vm.tracks) {
        let itemTitle = track.title
        if (track.title == "") {
          itemTitle = track.fileName
        }
        this.playlist.push({title: itemTitle, artist: track.artist, path: track.location})
      }
    },
    nextTrack: function() {
      if (this.currentHowl) {
        this.currentHowl.stop()
      }
      this.currentIndex++
      if (this.currentIndex >= this.playlist.length) {
        this.currentIndex = 0
      }
      this.currentHowl = this.makeNewHowl(this.currentItem.path)
      if (this.playing) {
        this.currentHowl.play()
      }
    },
    playPause: function() {
      if (this.playing) {
        this.currentHowl.pause()
        this.playing = false
      } else {
        this.currentHowl.play()
        this.playing = true
      }
    },
    updateVolume (level) {
      Howler.volume(level)
    },
  },
}
</script>

<style>
.active {
  font-weight: bold;
}
</style>
