
var diary = new Vue({
  el: '#diaryView',
  data: {
    diaries: {}
  },
  created () {
   fetch ("/diary")
   .then(response => response.json())
   .then(json => {
     this.diaries = json
   })
  },
  methods: {
    refresh: function () {
        console.log("Refresing requested.")
        fetch ("/diary?s_date=" + this.diaries.StartDate + "&e_date=" + this.diaries.EndDate)
        .then(response => response.json())
        .then(json => {
          this.diaries = json
        })
    },
    isMonday: function (str) {
      return new Date(str).getDay() === 0
    },
    isUpdated: function (int){
      return int > 0
    },
    append: function (){
      this.diaries.List.push({
        Id: 0,
        Content: null
      })
    },
    post: function (diary){
        fetch("/diary", {
        method: "POST",
        cache: "no-cache",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(diary),
      })
      .then(response => response.json())
      .then(json => {
        if (diary.Id == 0) {
          dialy.Id == json
        } else {
          diary.Updated = json
        }
      })
  }
  }
})

var transaction = new Vue({
  el: '#transaction',
  data: {
    transactions: {}
  },
  created () {
    fetch("/transaction")
    .then(response => response.json())
    .then(json => {
      this.transactions = json
    })
  },
  methods: {
    refresh: function () {
          fetch ("/transaction?s_date=" + this.transactions.StartDate + "&e_date=" + this.transactions.EndDate)
          .then(response => response.json())
          .then(json => {
            this.transactions = json
          })
    }
  }
})
