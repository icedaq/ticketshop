<template>
  <div class="hello">
    <h1>Get your ticket!</h1>

    <div class="loading" v-if="loading">
      Loading...
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>
    <div id="shows">
      <span>Name: {{ data[0].name }}</span>
      <span>Start time: {{ data[0].starttime }}</span>
      <span>End time: {{ data[0].endtime }}</span>
      <span>Price: {{ data[0].price }} CHF</span>
      <span>Max tickets: {{ data[0].maxtickets }}</span>
      <span><router-link :to="`/buy/${data[0].id}`" tag="button">BUY</router-link></span>
    </div>

  </div>
</template>

<script>
export default {
  name: 'Tickets',
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
      this.axios.get(this.$APIServerURL + '/tickets/' + showId).then((response) => {
        this.loading = false
        if (response.status !== 200) {
          this.error = 'Could not load data.'
        } else {
          this.data = response.data
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
