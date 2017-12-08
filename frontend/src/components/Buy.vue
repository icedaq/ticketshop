<template>
  <div class="hello">
    <h1>Purchase</h1>

    <div class="loading" v-if="loading">
      Purchase is being processed...
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>
    <div id="buy">
      {{ this.data }}
    </div>

  </div>
</template>

<script>
export default {
  name: 'Buy',
  data () {
    return {
      data: null,
      error: null
    }
  },
  created () {
    // fetch the data when the view is created and the data is
    // already being observed
    this.fetchData()
  },
  watch: {
    // call again the method if the route changes
    '$route': 'fetchData'
  },
  methods: {
    fetchData () {
      this.error = this.data = null
      this.loading = true

      var showId = this.$route.params.id
      this.axios.get(this.$APIServerURL + '/buy/' + showId).then((response) => {
        this.loading = false
        if (response.status !== 200) {
          this.error = 'Could not load data.'
        } else if (response.data.toString() === 'true') {
          this.data = 'Purchase successful!'
        } else if (response.data.toString() === 'false') {
          this.data = 'Purchase failed! Tickets sold out!'
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
  margin: 0;
  padding: 0;
}

a {
  color: #42b983;
}

#tickets {
  width: 200px;
  margin: 0 auto;
}

span {
  display: block;
}

</style>
