/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nutil

import "github.com/google/uuid"

func GetGUUID() string {
	return uuid.New().String()
}
