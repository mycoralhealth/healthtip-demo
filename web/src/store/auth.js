import User from '@/models/user'

import * as MutationTypes from './mutation-types'

const state = {
  user: User.from(localStorage.result)
}

const mutations = {
  [MutationTypes.LOGIN] (state) {
    state.user = User.from(localStorage.result)
  },
  [MutationTypes.LOGOUT] (state) {
    state.user = null
    delete localStorage.result
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
