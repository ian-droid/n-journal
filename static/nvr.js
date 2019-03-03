
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
          diary.Id =json
          diary.Updated = json
        }
        diary.Updated = json
      })
  }
  }
})

var transaction = new Vue({
  el: '#transactionView',
  data: {
    transactions: {},
    Currencies: {},
    Banks: {},
    Payments: {}
  },
  created () {
    fetch("/transaction")
    .then(response => response.json())
    .then(json => {
      this.transactions = json
    })
    fetch("/fininfo")
    .then(response => response.json())
    .then(json => {
      this.Currencies = json.Currencies
      this.Banks = json.Banks
      this.Payments = json.Payments
    })
  },
  methods: {
    refresh: function () {
          fetch ("/transaction?s_date=" + this.transactions.StartDate + "&e_date=" + this.transactions.EndDate)
          .then(response => response.json())
          .then(json => {
            this.transactions = json
          })
    },
    add: function (){
      this.transactions.List.push({
        Id: 0,
        Direction: "Pay",
        Currency: this.Currencies[0].Id,
        Bank: this.Banks[0].Id,
        Payment: this.Payments[0].Id
      })
    },
    post: function (transaction) {
      if (transaction.Direction === "Pay") {
        transaction.Pay = true
        transaction.Income = false
      } else if (transaction.Direction === "Income") {
        transaction.Pay = false
        transaction.Income = true
      }
      transaction.CurrencyPrefix = this.Currencies.find(c => c.Id === transaction.Currency).Prefix
      transaction.BankName = this.Banks.find(b => b.Id === transaction.Bank).Name
      transaction.PaymentName = this.Payments.find(p => p.Id === transaction.Payment).Name
      fetch("/transaction", {
      method: "POST",
      cache: "no-cache",
      headers: {
          "Content-Type": "application/json",
      },
      body: JSON.stringify(transaction),
    })
    .then(response => response.json())
    .then(json => {
      transaction.Id = json
      transaction.Updated = json
    })
    }
  }
})
