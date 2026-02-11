
const calculateDelay = (retryCount) => {
    return Math.pow(2, retryCount) * 1000;
};

