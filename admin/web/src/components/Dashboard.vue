<template>
    <div>
    <v-card v-for="(app, index) in apps" :key="app.Name">
      <v-card-row class="green darken-1">
        <v-card-title>
            <span class="white--text">{{ app.Name }}</span>
        </v-card-title>
      </v-card-row>
      <v-card-text>
        <v-card-row height="75px">
          <div>
          </div>
        </v-card-row>
      </v-card-text>
      <v-card-row actions>
        <v-btn flat class="green--text darken-1">Setting</v-btn>
      </v-card-row>
    </v-card>
    </div>
</template>

<script>

    import config from '../config'
    import app from '../store/app'

    export default {
        data () {
            return {
                itemGroup: config.menu,
                apps: []
            }
        },
        created () {
            this.bindProjects()
        },

        methods: {
            bindProjects () {
                let my = this
                app.list(function (response) {
                    my.apps = response.data
                })
            }
        }
    }
</script>
