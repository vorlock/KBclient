// @flow
import * as Types from '../../../constants/types/fs'
import * as Constants from '../../../constants/fs'
import {isMobile, isIOS} from '../../../constants/platform'

export type Layout = {
  copyPath: boolean,
  delete: boolean,
  download: boolean,
  ignoreTlf: boolean,
  moveOrCopy: boolean,
  saveMedia: boolean,
  showInSystemFileManager: boolean,
  // if multiple share items exist, they go into 2nd level menu - share
  sendAttachmentToChat: boolean,
  sendLinkToChat: boolean,
  sendToOtherApp: boolean,
  share: boolean,
}

const empty = {
  copyPath: false,
  delete: false,
  download: false,
  ignoreTlf: false,
  moveOrCopy: false,
  saveMedia: false,
  showInSystemFileManager: false,
  // share items
  // eslint-disable-next-line sort-keys
  sendAttachmentToChat: false,
  sendLinkToChat: false,
  sendToOtherApp: false,
  share: false,
}

const getRawLayout = (path: Types.Path, pathItem: Types.PathItem): Layout => {
  const level = Types.getPathLevel(path)
  switch (level) {
    case 0:
    case 1:
      // should never happen
      return empty
    case 2: // private/public/team
      return {
        ...empty,
        copyPath: true,
        showInSystemFileManager: !isMobile,
      }
    case 3: // tlf
      return {
        ...empty,
        copyPath: true,
        ignoreTlf: true,
        sendLinkToChat: false, // TODO enable mobile // desktop uses separate button
        showInSystemFileManager: !isMobile,
      }
    default:
      // inside tlf
      return {
        ...empty,
        copyPath: true,
        delete: pathItem.type === 'file',
        download: pathItem.type === 'file' && !isIOS,
        moveOrCopy: true,
        saveMedia: isMobile && pathItem.type === 'file' && Constants.isMedia(pathItem),
        showInSystemFileManager: !isMobile,
        // share menu items
        // eslint-disable-next-line sort-keys
        sendAttachmentToChat: false, // TODO enable mobile pathItem.type === 'file' && isMobile, // desktop uses separate button
        sendLinkToChat: false, // TODO enable mobile // desktop uses separate button
        sendToOtherApp: pathItem.type === 'file' && isMobile,
      }
  }
}

const totalShare = layout =>
  (layout.sendAttachmentToChat ? 1 : 0) + (layout.sendLinkToChat ? 1 : 0) + (layout.sendToOtherApp ? 1 : 0)

const consolidateShares = (layout: Layout): Layout =>
  totalShare(layout) > 1
    ? {
        ...layout,
        sendAttachmentToChat: false,
        sendLinkToChat: false,
        sendToOtherApp: false,
        share: true,
      }
    : layout

const filterForOnlyShares = (layout: Layout): Layout => ({
  ...empty,
  sendAttachmentToChat: layout.sendAttachmentToChat,
  sendLinkToChat: layout.sendLinkToChat,
  sendToOtherApp: layout.sendToOtherApp,
})

export const getRootLayout = (path: Types.Path, pathItem: Types.PathItem): Layout =>
  consolidateShares(getRawLayout(path, pathItem))

export const getShareLayout = (path: Types.Path, pathItem: Types.PathItem): Layout =>
  filterForOnlyShares(getRawLayout(path, pathItem))
