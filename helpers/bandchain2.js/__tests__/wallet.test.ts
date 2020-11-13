import { PrivateKey } from '../src/wallet'

const TEST_MNEMONIC =
  'coach pond canoe lake solution empty vacuum term pave toe burst top violin purpose umbrella color disease thrive diamond found track need filter wait'

describe('Private Key', () => {
  it("fromMnemonic with m/44'/494'/0'/0/0", () => {
    const privkey = PrivateKey.fromMnemonic(TEST_MNEMONIC)

    expect(privkey.toHex()).toEqual(
      '2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d',
    )
  })

  it("fromMnemonic with m/44'/494'/0'/0/1", () => {
    const privkey = PrivateKey.fromMnemonic(TEST_MNEMONIC, "m/44'/494'/0'/0/1")

    expect(privkey.toHex()).toEqual(
      '987af53f91a09926274e5a2ef86223356112056f61b35a57df345d7b14176bb3',
    )
  })

  it("generate with m/44'/494'/0'/0/1", () => {
    const [mnemonic, privkey] = PrivateKey.generate("m/44'/494'/0'/0/5")

    expect(
      PrivateKey.fromMnemonic(mnemonic, "m/44'/494'/0'/0/5").toHex(),
    ).toEqual(privkey.toHex())
  })

  it('fromHex', () => {
    expect(
      PrivateKey.fromHex(
        '2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d',
      ).toHex(),
    ).toEqual(
      '2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d',
    )
  })

  it('to public key', () => {
    expect(PrivateKey.fromMnemonic(TEST_MNEMONIC).toPubkey().toHex()).toEqual(
      '0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d',
    )
  })

  it('sign by private key', () => {
    const privkey = PrivateKey.fromMnemonic(TEST_MNEMONIC)
    expect(
      privkey.sign(Buffer.from('test msg', 'utf-8')).toString('hex'),
    ).toEqual(
      '42a1e41012155ae2daa9b9a2e038f76463da4662564b4989f236ecb4d97f592c1190d92319363e2d1eb78fb98f0dac30c5e2a850f45bb4c44f1f6203ebe6efbe',
    )
  })
})

describe('Public Key', () => {
  it('verify signature', () => {
    const privkey = PrivateKey.fromMnemonic(TEST_MNEMONIC)
    const pubkey = privkey.toPubkey()
    const msg = Buffer.from('test msg', 'utf-8')
    const signature = privkey.sign(msg)

    expect(pubkey.verify(msg, signature)).toBe(true)
  })
})