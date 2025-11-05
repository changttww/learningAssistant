import { computed, reactive } from "vue";
import { getUserProfile, getUserStudyStats } from "@/api/modules/user";
import { getUserInfo as getStoredUserInfo, setUserInfo as storeUserInfo } from "@/utils/auth";

export const DEFAULT_USER_ID = 1;

const storedProfile = getStoredUserInfo();
const hasHydratedProfile = storedProfile && storedProfile.basic_info;

const state = reactive({
  profile: storedProfile,
  profileLoaded: !!hasHydratedProfile,
  profileLoading: false,
  profileError: "",
  studyStats: null,
  studyStatsLoaded: false,
  studyStatsLoading: false,
  studyStatsError: "",
});

let pendingProfilePromise = null;
let pendingStatsPromise = null;

function assignProfile(profile) {
  state.profile = profile;
  state.profileLoaded = !!profile?.basic_info;
  storeUserInfo(profile || null);
}

export function useCurrentUser() {
  async function loadCurrentUser(userId, options = {}) {
    const targetId =
      typeof userId === "number"
        ? userId
        : state.profile?.id || DEFAULT_USER_ID;
    const { force = false } = options;

    if (state.profileLoaded && !force && pendingProfilePromise === null) {
      return state.profile;
    }

    if (pendingProfilePromise && !force) {
      return pendingProfilePromise;
    }

    state.profileLoading = true;
    state.profileError = "";

    pendingProfilePromise = getUserProfile(targetId)
      .then((res) => {
        assignProfile(res.data || null);
        return state.profile;
      })
      .catch((error) => {
        state.profileError = error?.message || "获取用户信息失败";
        throw error;
      })
      .finally(() => {
        state.profileLoading = false;
        pendingProfilePromise = null;
      });

    return pendingProfilePromise;
  }

  async function loadStudyStats(userId, options = {}) {
    const targetId =
      typeof userId === "number"
        ? userId
        : state.profile?.id || DEFAULT_USER_ID;
    const { force = false } = options;

    if (state.studyStatsLoaded && !force && pendingStatsPromise === null) {
      return state.studyStats;
    }

    if (pendingStatsPromise && !force) {
      return pendingStatsPromise;
    }

    state.studyStatsLoading = true;
    state.studyStatsError = "";

    pendingStatsPromise = getUserStudyStats(targetId)
      .then((res) => {
        state.studyStats = res.data || null;
        state.studyStatsLoaded = !!res.data;
        return state.studyStats;
      })
      .catch((error) => {
        state.studyStatsError = error?.message || "获取学习统计失败";
        throw error;
      })
      .finally(() => {
        state.studyStatsLoading = false;
        pendingStatsPromise = null;
      });

    return pendingStatsPromise;
  }

  function setCurrentUser(profile) {
    assignProfile(profile);
  }

  function clearCurrentUser() {
    state.profile = null;
    state.profileLoaded = false;
    state.profileLoading = false;
    state.profileError = "";
    state.studyStats = null;
    state.studyStatsLoaded = false;
    state.studyStatsLoading = false;
    state.studyStatsError = "";
    storeUserInfo(null);
  }

  return {
    profile: computed(() => state.profile),
    profileLoaded: computed(() => state.profileLoaded),
    profileLoading: computed(() => state.profileLoading),
    profileError: computed(() => state.profileError),
    studyStats: computed(() => state.studyStats),
    studyStatsLoaded: computed(() => state.studyStatsLoaded),
    studyStatsLoading: computed(() => state.studyStatsLoading),
    studyStatsError: computed(() => state.studyStatsError),
    loadCurrentUser,
    loadStudyStats,
    setCurrentUser,
    clearCurrentUser,
  };
}
