var store = {
    state: {
        attendees: [],
        count: 3,
    }
}

new Vue({
    el: '#navcard',
    data: {
        iscard: true,
        shared: store.state,
    },
    methods: {
        fetch: function () {
            axios.get('/event')
                .then(resp => {
                    resp.data.forEach(e => {
                        e.member.won = false
                        e.member.picked = false
                        this.shared.attendees.push(e.member)
                    })
                    this.iscard = false
                })
                .catch(error => {
                    alert(error.response.data.message)
                })
        },
        lotto: function () {
            var self = this
            var len = this.shared.attendees.length
            var cnt = parseInt(Math.random() * len * 20)
            while (this.shared.count) {
                var i = -1
                while (true) {
                    i++
                    setTimeout(function (i) {
                        self.shared.attendees[i % len].picked = true;
                        if (i > 0) {
                            self.shared.attendees[(i - 1) % len].picked = false;
                        }
                    }, 12, i)
                    if (i < cnt - 1) {
                        continue
                    }

                    cnt = parseInt(Math.random() * len * 20)
                    if (this.shared.attendees[i % len].won) {
                        continue
                    }
                    setTimeout(function (i) {
                        self.shared.attendees[i % len].won = true
                    }, 10, i)
                    break
                }
                this.shared.count -= 1
            }
        }
    }
})

new Vue({
    el: '#members',
    data: {
        shared: store.state,
    },
    filters: {
        truncate: function (str, limit) {
            if (str.length > limit) {
                str = str.substring(0, limit - 3) + '...'
            }
            return str
        }
    },
    methods: {
        exclude: function (index) {
            if (this.shared.attendees[index].won) {
                this.shared.count += 1
            }
            this.shared.attendees.splice(index, 1)
        }
    }
})