import { useState, useEffect } from 'react'

/**
 * Custom hook to detect when a key is pressed
 * @param {string} targetKey - The key to listen for
 * @returns {boolean} Whether the key is currently pressed
 */
export const useKeyPress = (targetKey) => {
  const [keyPressed, setKeyPressed] = useState(false)

  useEffect(() => {
    const downHandler = ({ key }) => {
      if (key === targetKey) {
        setKeyPressed(true)
      }
    }

    const upHandler = ({ key }) => {
      if (key === targetKey) {
        setKeyPressed(false)
      }
    }

    window.addEventListener('keydown', downHandler)
    window.addEventListener('keyup', upHandler)

    return () => {
      window.removeEventListener('keydown', downHandler)
      window.removeEventListener('keyup', upHandler)
    }
  }, [targetKey])

  return keyPressed
}

/**
 * Custom hook to detect keyboard shortcuts
 * @param {string[]} keys - Array of keys to listen for (e.g., ['Control', 's'])
 * @param {function} callback - Callback function when keys are pressed
 * @param {object} options - Options for the hook
 */
export const useKeyboardShortcut = (keys, callback, options = {}) => {
  const { preventDefault = true, enabled = true } = options
  const [keysPressed, setKeysPressed] = useState(new Set())

  useEffect(() => {
    if (!enabled) return

    const downHandler = (event) => {
      const key = event.key
      setKeysPressed((prev) => new Set(prev).add(key))

      // Check if all keys in the shortcut are pressed
      const allKeysPressed = keys.every((k) => {
        if (k === key) return true
        return keysPressed.has(k) || k === key
      })

      if (allKeysPressed) {
        if (preventDefault) {
          event.preventDefault()
        }
        callback(event)
      }
    }

    const upHandler = ({ key }) => {
      setKeysPressed((prev) => {
        const next = new Set(prev)
        next.delete(key)
        return next
      })
    }

    window.addEventListener('keydown', downHandler)
    window.addEventListener('keyup', upHandler)

    return () => {
      window.removeEventListener('keydown', downHandler)
      window.removeEventListener('keyup', upHandler)
    }
  }, [keys, callback, preventDefault, enabled, keysPressed])
}

/**
 * Hook for Escape key
 */
export const useEscapeKey = (callback) => {
  useEffect(() => {
    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        callback(event)
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [callback])
}

