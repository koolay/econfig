<template>
    <div>
    <v-card>
      <v-card-row class="green darken-1">
        <v-card-title>
          <span class="white--text">Marriot Rewards</span>
        </v-card-title>
      </v-card-row>
      <v-card-text>
        <v-card-row height="75px">
          <div>
            <div>Membership Number</div><strong>113241423</strong>
          </div>
        </v-card-row>
      </v-card-text>
      <v-card-row actions>
        <v-btn flat class="green--text darken-1">View Email</v-btn>
      </v-card-row>
    </v-card>

    <v-card>
      <v-card-row class="green darken-1">
        <v-card-title>
          <span class="white--text">Marriot Rewards</span>
        </v-card-title>
      </v-card-row>
      <v-card-text>
        <v-card-row height="75px">
          <div>
            <div>Membership Number</div><strong>113241423</strong>
          </div>
        </v-card-row>
      </v-card-text>
      <v-card-row actions>
        <v-btn flat class="green--text darken-1">View Email</v-btn>
      </v-card-row>
    </v-card>
    <v-card>
      <v-card-row class="green darken-1">
        <v-card-title>
          <span class="white--text">Marriot Rewards</span>
        </v-card-title>
      </v-card-row>
      <v-card-text>
        <v-card-row height="75px">
          <div>
            <div>Membership Number</div><strong>113241423</strong>
          </div>
        </v-card-row>
      </v-card-text>
      <v-card-row actions>
        <v-btn flat class="green--text darken-1">View Email</v-btn>
      </v-card-row>
    </v-card>
    </div>
</template>

<script>

    import config from '../config'
    import project from '../store/project'

    export default {
        data () {
            return {
                itemGroup: config.menu,
                projects: []
            }
        },
        created () {
            this.bindProjects()
        },

        methods: {
            bindProjects () {
                project.list(function (data) {
                    this.projects = data
                })
            }
        }
    }
</script>
