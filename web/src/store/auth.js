import User from '@/models/user'

import * as MutationTypes from './mutation-types'

const state = {
  user: User.from(localStorage.result)
}

const mutations = {
  [MutationTypes.LOGIN] (state) {
  	console.log(localStorage.result)
    state.user = User.from(localStorage.result)
    console.log("state.user")
    console.log(state.user)
  },
  [MutationTypes.LOGOUT] (state) {
    state.user = null
  }
}

const getters = {
  currentUser (state) {
    return state.user
  }
}

const actions = {
  login ({ commit }) {
    commit(MutationTypes.LOGIN)
  },

  logout ({ commit }) {
    commit(MutationTypes.LOGOUT)
  }
}

export default {
  state,
  mutations,
  getters,
  actions
}
